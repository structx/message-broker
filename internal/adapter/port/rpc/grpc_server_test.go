package rpc_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/trevatk/go-pkg/adapter/logging"
	"github.com/trevatk/go-pkg/adapter/setup"
	pb "github.com/trevatk/go-pkg/proto/messaging/v1"
	"github.com/trevatk/go-pkg/util/decode"
	"github.com/trevatk/mora/internal/adapter/port/rpc"
	"github.com/trevatk/mora/internal/core/domain"
)

func init() {
	_ = os.Setenv("ROOT_CONFIG", "./testfiles/test_config.hcl")
}

type GRPCServerSuite struct {
	suite.Suite
	s pb.MessagingServiceV1Server
}

func (suite *GRPCServerSuite) SetupTest() {

	assert := suite.Assert()

	logger, err := logging.New(nil)
	assert.NoError(err)

	cfg := setup.New()
	assert.NoError(decode.ConfigFromEnv(cfg))

	mockRaft := domain.NewMockRaft(suite.T())

	suite.s = rpc.NewGRPCServer(logger, cfg, mockRaft)
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
