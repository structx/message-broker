package rpc_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/trevatk/go-pkg/logging"
	pb "github.com/trevatk/go-pkg/proto/messaging/v1"
	"github.com/trevatk/mora/internal/adapter/port/rpc"
	"github.com/trevatk/mora/internal/adapter/setup"
	"github.com/trevatk/mora/internal/core/domain"
)

func init() {
	_ = os.Setenv("SERVER_GRPC_PORT", "50051")
}

type GRPCServerSuite struct {
	suite.Suite
	s pb.MessagingServiceV1Server
}

func (suite *GRPCServerSuite) SetupTest() {

	assert := suite.Assert()
	ctx := context.TODO()

	logger, err := logging.NewLoggerFromEnv()
	assert.NoError(err)

	cfg := setup.NewConfig()
	assert.NoError(setup.ProcessConfigWithEnv(ctx, cfg))

	mockInterceptor := domain.NewMockAuthenticatorInterceptor(suite.T())

	mockRaft := domain.NewMockRaft(suite.T())

	suite.s = rpc.NewGRPCServer(logger, cfg, mockInterceptor, mockRaft)
}

func (suite *GRPCServerSuite) TestPublish() {

	ctx := context.TODO()

	assert := suite.Assert()

	testcases := []struct {
		expected error
		in       *pb.Envelope
	}{
		{
			in: &pb.Envelope{
				Topic:   "test.publish",
				Payload: []byte("hello world"),
			},
			expected: nil,
		},
		{
			// invalid topic name
			in: &pb.Envelope{
				Topic: "unittest",
			},
			expected: status.Error(codes.InvalidArgument, "invalid topic parameter"),
		},
	}

	for _, tc := range testcases {

		_, err := suite.s.Publish(ctx, tc.in)
		assert.Equal(tc.expected, err)

	}
}

func (suite *GRPCServerSuite) TestSubscribe() {}

func (suite *GRPCServerSuite) TestRequestResponse() {}

func TestGRPCServerSuite(t *testing.T) {
	suite.Run(t, new(GRPCServerSuite))
}
