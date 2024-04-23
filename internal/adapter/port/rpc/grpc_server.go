// Package rpc implementation of rpc using gRPC
package rpc

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/trevatk/go-pkg/proto/messaging/v1"
	"github.com/trevatk/go-pkg/structs/dht"

	"github.com/trevatk/mora/internal/adapter/setup"
	"github.com/trevatk/mora/internal/core/domain"
)

// GRPCServer implementation of messaging service gRPC interface
type GRPCServer struct {
	// server interface compliance
	pb.UnimplementedMessagingServiceV1Server

	log       *zap.SugaredLogger
	dht       *dht.DHT
	mtx       sync.RWMutex
	subs      map[string][]pb.MessagingServiceV1_SubscribeServer
	envelopes map[string][]*pb.Envelope
	r         domain.Raft

	port int
	s    *grpc.Server
}

// NewGRPCServer return new gRPC server class
func NewGRPCServer(logger *zap.Logger, cfg *setup.Config, auth domain.AuthenticatorInterceptor, raft domain.Raft) *GRPCServer {
	return &GRPCServer{
		log:       logger.Sugar().Named("grpc_server"),
		mtx:       sync.RWMutex{},
		subs:      make(map[string][]pb.MessagingServiceV1_SubscribeServer),
		envelopes: make(map[string][]*pb.Envelope),
		dht:       dht.NewDHT(cfg.Server.Addr),
		port:      cfg.Server.GRPCPort,
		r:         raft,
		s:         grpc.NewServer(grpc.UnaryInterceptor(auth.UnaryInterceptor), grpc.StreamInterceptor(auth.StreamInterceptor)),
	}
}

// Publish add message to chain and send message to subscribers
func (g *GRPCServer) Publish(ctx context.Context, in *pb.Envelope) (*pb.Stub, error) {

	// notify all nodes in consensus
	g.r.Apply(ctx, nil)

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
func (g *GRPCServer) RequestResponse(_ context.Context, _ *pb.Envelope) (*pb.Envelope, error) {
	// TODO:
	// implement handler
	return nil, nil
}

// Start gRPC server
func (g *GRPCServer) Start() error {

	pb.RegisterMessagingServiceV1Server(g.s, g)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		return fmt.Errorf("failed to create network listener %v", err)
	}

	go func() {
		if err := g.s.Serve(listener); err != nil {
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
