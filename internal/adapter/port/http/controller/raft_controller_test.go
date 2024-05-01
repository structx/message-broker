package controller_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/structx/go-pkg/adapter/logging"

	"github.com/structx/message-broker/internal/adapter/port/http/controller"
	"github.com/structx/message-broker/internal/core/domain"
)

func init() {
	_ = os.Setenv("LOG_LEVEL", "DEBUG")
	_ = os.Setenv("LOG_PATH", "./testfiles/raft_test.log")
}

type RaftControllerSuite struct {
	suite.Suite
	handler http.Handler
}

func (suite *RaftControllerSuite) SetupTest() {

	assert := suite.Assert()

	logger, err := logging.New(nil)
	assert.NoError(err)

	mockRaft := domain.NewMockRaft(suite.T())
	mockRaft.EXPECT().Join(mock.AnythingOfType("context.Context"), mock.AnythingOfType("*domain.NewMember")).Return(&domain.Member{
		Index:    1,
		ServerID: "12345678901234567",
		Addr:     "127.0.0.1:8333",
	}, nil).Once()

	rc := controller.NewRaft(logger, mockRaft)
	suite.handler = rc.RegisterRoutesV1P()
}

func (suite *RaftControllerSuite) TestJoin() {

	assert := suite.Assert()

	tt := []struct {
		expected int
		params   controller.JoinRequestParams
	}{
		{
			expected: http.StatusCreated,
			params: controller.JoinRequestParams{
				NewMember: &controller.NewMember{
					ServerID: "12345678901234567",
					Address:  "127.0.0.1:8333",
				},
			},
		},
	}

	for _, tc := range tt {

		rr := httptest.NewRecorder()

		requestbytes, err := json.Marshal(tc.params)
		assert.NoError(err)

		request, err := http.NewRequest(http.MethodPut, "/raft/join", bytes.NewBuffer(requestbytes))
		assert.NoError(err)

		suite.handler.ServeHTTP(rr, request)

		assert.Equal(tc.expected, rr.Code)

		if tc.expected == http.StatusCreated {
			body, err := io.ReadAll(rr.Body)
			assert.NoError(err)

			var response controller.JoinResponse
			err = json.Unmarshal(body, &response)
			assert.NoError(err)

			assert.Equal(tc.params.NewMember.Address, response.Member.Address)
			assert.Equal(tc.params.NewMember.ServerID, response.Member.ServerID)
		}
	}
}

func TestRaftControllerSuite(t *testing.T) {
	suite.Run(t, new(RaftControllerSuite))
}
