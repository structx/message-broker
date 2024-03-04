// Package rpc implementation of rpc using gRPC
package rpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/trevatk/block-broker/internal/adapter/port/rpc/proto/messaging/v1"
	"github.com/trevatk/block-broker/internal/adapter/setup"
	"github.com/trevatk/block-broker/internal/core/domain"
)

// GRPCServer implementation of messaging service gRPC interface
type GRPCServer struct {
	// server interface compliance
	pb.UnimplementedMessagingServiceServer

	m domain.Messenger

	log *zap.SugaredLogger

	mtx       sync.Mutex
	subs      sync.Map
	envelopes sync.Map

	port int
	s    *grpc.Server
}

// NewGRPCServer return new gRPC server class
func NewGRPCServer(logger *zap.Logger, cfg *setup.Config, messenger domain.Messenger) *GRPCServer {
	return &GRPCServer{
		log:       logger.Sugar().Named("grpc_server"),
		m:         messenger,
		mtx:       sync.Mutex{},
		subs:      sync.Map{},
		envelopes: sync.Map{},
		port:      cfg.Server.GRPCPort,
	}
}

// Publish add message to chain and send message to subscribers
func (g *GRPCServer) Publish(_ context.Context, in *pb.Envelope) (*pb.Stub, error) {

	topic := in.GetTopic()
	payload := in.GetPayload()

	if topic == "" || !strings.Contains(topic, ".") {
		return nil, status.Error(codes.InvalidArgument, "invalid topic parameter")
	}

	// md, ok := metadata.FromIncomingContext(ctx)
	// if !ok {
	// 	g.log.Error("invalid metadata provided")
	// 	return nil, status.Error(codes.InvalidArgument, "invalid metadata provided")
	// }

	msg, err := g.m.Create(&domain.NewMessage{
		Topic:     topic,
		Payload:   payload,
		Publisher: "test",
	})
	if err != nil {
		g.log.Errorf("g.m.Create: %v", err)
		return nil, status.Error(codes.Internal, "unable to publish message")
	}

	if v, ok := g.envelopes.LoadOrStore(topic, []*pb.Envelope{}); !ok {
		// existing map found
		if envs, ok := v.([]*pb.Envelope); ok {
			envs = append(envs, in)
			g.envelopes.Store(topic, envs)
		}
	} else {
		// empty map loaded
		// load first envelope
		g.envelopes.Store(topic, []*pb.Envelope{in})
	}

	return &pb.Stub{
		EnvelopeId: msg.ID,
	}, nil
}

// Subscribe store subscription in memory
func (g *GRPCServer) Subscribe(in *pb.Subscription, stream pb.MessagingService_SubscribeServer) error {

	ctx := stream.Context()

	topic := in.GetTopic()

	if err := isValidTopic(topic); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if v, ok := g.subs.LoadOrStore(topic, []pb.MessagingService_SubscribeServer{}); !ok {
		if subs, ok := v.([]pb.MessagingService_SubscribeServer); ok {
			subs = append(subs, stream)
			g.subs.Store(topic, subs)
		}
	} else {
		g.subs.Store(topic, []pb.MessagingService_SubscribeServer{stream})
	}

OUTER:
	for {

		select {
		case <-ctx.Done():
			// TODO:
			// remove subscriber
			return nil
		default:

			v, ok := g.envelopes.Load(topic)
			if !ok {
				continue OUTER
			}

			envs, ok := v.([]*pb.Envelope)
			if !ok {
				continue OUTER
			}

			for _, msg := range envs {

				v, ok := g.subs.Load(topic)
				if !ok {
					continue OUTER
				}

				subs, ok := v.([]pb.MessagingService_SubscribeServer)
				if !ok {
					continue OUTER
				}

			SUBSCRIPTIONS:
				for _, sub := range subs {
					err := sub.Send(msg)
					if err != nil {
						continue SUBSCRIPTIONS
					}
				}
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

	g.s = grpc.NewServer()
	pb.RegisterMessagingServiceServer(g.s, g)

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

func isValidTopic(topic string) error {

	if len(topic) < 1 {
		return errors.New("invalid topic length")
	} else if !strings.Contains(topic, ".") {
		return errors.New("topic does not contain noun followed by verb")
	}

	return nil
}
