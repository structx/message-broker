package controller_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/trevatk/block-broker/internal/adapter/logging"
	"github.com/trevatk/block-broker/internal/adapter/port/http/router"
	"github.com/trevatk/block-broker/internal/core/domain"
)

type MessageControllerSuite struct {
	suite.Suite
	handler http.Handler
}

func (suite *MessageControllerSuite) SetupTest() {

	assert := suite.Assert()

	logger, err := logging.NewLogger()
	assert.NoError(err)

	mockMessenger := domain.NewMockMessenger(suite.T())
	mockMessenger.EXPECT().Read(mock.Anything).Return(&domain.Message{}, nil).Once()

	suite.handler = router.NewRouter(logger, mockMessenger)
}

func (suite *MessageControllerSuite) TestFetchMessage() {

	assert := suite.Assert()

	testcases := []struct {
		hash     string
		expected int
	}{
		{
			hash:     "0cf6b052518a08ed5299676530c3d57f2c94dbe9ec26184b4b9d3baf",
			expected: http.StatusAccepted,
		},
		{
			hash:     "",
			expected: http.StatusNotFound,
		},
	}

	for _, tt := range testcases {

		rr := httptest.NewRecorder()

		endpoint := fmt.Sprintf("/api/v1/messages/%s", tt.hash)
		request, err := http.NewRequest(http.MethodGet, endpoint, nil)
		assert.NoError(err)

		suite.handler.ServeHTTP(rr, request)

		assert.Equal(tt.expected, rr.Code)
	}
}

func (suite *MessageControllerSuite) TestListMessages() {}

func (suite *MessageControllerSuite) TestListTopics() {}

func (suite *MessageControllerSuite) TestListByTopic() {}

func TestMessageControllerSuite(t *testing.T) {
	suite.Run(t, new(MessageControllerSuite))
}
