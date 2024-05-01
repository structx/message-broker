// Package rpc implementation of rpc using gRPC
package rpc

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/leaderhealth"
	"github.com/Jille/raftadmin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/structx/go-pkg/proto/messaging/v1"
	"github.com/structx/go-pkg/structs/dht"

	pkgdomain "github.com/structx/go-pkg/domain"
	"github.com/structx/message-broker/internal/core/domain"
)

// GRPCServer implementation of messaging service gRPC interface
type GRPCServer struct {
	// server interface compliance
	pb.UnimplementedMessagingServiceV1Server

	log       *zap.SugaredLogger
	dht       map[string]*dht.DHT
	mtx       sync.RWMutex
	subs      map[string][]pb.MessagingServiceV1_SubscribeServer
	envelopes map[string][]*pb.Envelope
	r         domain.Raft

	port int
	s    *grpc.Server
}

// NewGRPCServer return new gRPC server class
func NewGRPCServer(logger *zap.Logger, cfg pkgdomain.Config, raft domain.Raft) *GRPCServer {
	scfg := cfg.GetServer()
	return &GRPCServer{
		log:       logger.Sugar().Named("GrpcServer"),
		mtx:       sync.RWMutex{},
		subs:      make(map[string][]pb.MessagingServiceV1_SubscribeServer),
		envelopes: make(map[string][]*pb.Envelope),
		dht:       make(map[string]*dht.DHT),
		port:      scfg.Ports.GRPC,
		r:         raft,
		s:         grpc.NewServer(),
	}
}

// Publish add message to chain and send message to subscribers
func (g *GRPCServer) Publish(_ context.Context, in *pb.Envelope) (*pb.Stub, error) {

	// notify all local services subscribed
	g.mtx.Lock()
	defer g.mtx.Unlock()
	if _, ok := g.envelopes[in.Topic]; !ok {
		g.envelopes[in.Topic] = []*pb.Envelope{in}
	} else {
		g.envelopes[in.Topic] = append(g.envelopes[in.Topic], in)
	}

	return &pb.Stub{}, nil
}

// Subscribe store subscription in memory
func (g *GRPCServer) Subscribe(in *pb.Subscription, stream pb.MessagingServiceV1_SubscribeServer) error {

	ctx := stream.Context()

	p, ok := peer.FromContext(ctx)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "missing peer data")
	}

	n := dht.NewNode(p.Addr.String())
	err := g.dht[in.Topic].AddNode(n)
	if err != nil {
		g.log.Errorf("failed to add node to DHT %v", err)
		return status.Errorf(codes.Internal, "unable to add node to DHT")
	}

	g.mtx.Lock()
	if _, ok := g.subs[in.Topic]; !ok {
		g.subs[in.Topic] = []pb.MessagingServiceV1_SubscribeServer{stream}
	} else {
		g.subs[in.Topic] = append(g.subs[in.Topic], stream)
	}
	g.mtx.Unlock()

	for {
		select {
		case <-ctx.Done():
			// client disconnect
			// remove node from subscribers list
			for i, s := range g.subs[in.Topic] {
				if s == stream {
					g.subs[in.Topic] = removeSubscriber(g.subs[in.Topic], i)
				}
			}
			return nil
		default:
			for i, m := range g.envelopes[in.Topic] {
			SUBSCRIBERS:
				for j, s := range g.subs[in.Topic] {
					err := s.SendMsg(m)
					if err != nil {
						// removed errored node
						g.log.Errorf("failed to send message %v", err)
						g.subs[in.Topic] = removeSubscriber(g.subs[in.Topic], j)
						continue SUBSCRIBERS
					}
				}
				g.envelopes[in.Topic] = removeMessage(g.envelopes[in.Topic], i)
			}

			time.Sleep(time.Millisecond * 200)
		}

	}
}

// RequestResponse message handler
func (g *GRPCServer) RequestResponse(ctx context.Context, in *pb.Envelope) (*pb.Envelope, error) {

	// find closest node in topic
	dht := g.dht[in.Topic]
	n, err := dht.Put([]byte(in.Topic))
	if err != nil {
		g.log.Errorf("failed to resolve key %v", err)
		return nil, status.Errorf(codes.Internal, "no nodes exist")
	}

	// set s2s timeout
	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	// dial closest node
	conn, err := grpc.DialContext(timeout, n.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		g.log.Errorf("failed to dial addr %v", err)
		return nil, status.Errorf(codes.Internal, "unable to dial requested service")
	}
	defer func() { _ = conn.Close() }()

	// publish request with altered topic
	cli := pb.NewMessagingServiceV1Client(conn)
	result, err := cli.RequestResponse(ctx, in)
	if err != nil {
		g.log.Errorf("failed to send request/response to node %v", err)

	}

	return result, nil
}

// Start gRPC server
func (g *GRPCServer) Start(params *domain.GrpcStartParams) error {

	pb.RegisterMessagingServiceV1Server(g.s, g)
	params.TransportManager.Register(g.s)

	leaderhealth.Setup(
		params.Raft,
		g.s,
		[]string{"MessagingService"},
	)

	raftadmin.Register(g.s, params.Raft)
	reflection.Register(g.s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		return fmt.Errorf("failed to create network listener %v", err)
	}

	go func() {
		if err := g.s.Serve(lis); err != nil {
			g.log.Fatalf("unable to start gRPC server %v", err)
		}
	}()

	return nil
}

// Shutdown gRPC server
func (g *GRPCServer) Shutdown() {
	g.s.GracefulStop()
}

func removeMessage(s []*pb.Envelope, index int) []*pb.Envelope {
	if index < 0 || index >= len(s) {
		return s
	}
	return append(s[:index], s[index+1:]...)
}

func removeSubscriber(s []pb.MessagingServiceV1_SubscribeServer, index int) []pb.MessagingServiceV1_SubscribeServer {
	if index < 0 || index >= len(s) {
		return s
	}
	return append(s[:index], s[index+1:]...)
}
