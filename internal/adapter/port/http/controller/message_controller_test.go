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
	mockMessenger.On("Read", mock.Anything).Return(&domain.Message{}, nil).Maybe()
	mockMessenger.On("List", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return([]*domain.PartialMessage{}, nil).Maybe()
	mockMessenger.On("ListByTopic", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return([]*domain.PartialMessage{}, nil).Maybe()
	mockMessenger.On("ListTopics", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return([]string{}, nil).Maybe()

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

func (suite *MessageControllerSuite) TestListMessages() {

	assert := suite.Assert()

	testcases := []struct {
		expected      int
		limit, offset int
	}{
		{
			expected: http.StatusAccepted,
			limit:    10,
			offset:   0,
		},
	}

	for _, tt := range testcases {

		rr := httptest.NewRecorder()

		endpoint := fmt.Sprintf("/api/v1/messages?limit=%d&offset=%d", tt.limit, tt.offset)
		request, err := http.NewRequest(http.MethodGet, endpoint, nil)
		assert.NoError(err)

		suite.handler.ServeHTTP(rr, request)

		assert.Equal(tt.expected, rr.Code)
	}
}

func (suite *MessageControllerSuite) TestListTopics() {

	assert := suite.Assert()

	testcases := []struct {
		expected      int
		limit, offset int
	}{
		{
			expected: http.StatusAccepted,
			limit:    10,
			offset:   0,
		},
	}

	for _, tt := range testcases {

		rr := httptest.NewRecorder()

		endpoint := fmt.Sprintf("/api/v1/messages/topics?limit=%d&offset=%d", tt.limit, tt.offset)
		request, err := http.NewRequest(http.MethodGet, endpoint, nil)
		assert.NoError(err)

		suite.handler.ServeHTTP(rr, request)

		assert.Equal(tt.expected, rr.Code)
	}
}

func (suite *MessageControllerSuite) TestListByTopic() {

	assert := suite.Assert()

	testcases := []struct {
		expected      int
		limit, offset int
		topic         string
	}{
		{
			expected: http.StatusAccepted,
			limit:    10,
			offset:   0,
			topic:    "unit.test",
		},
	}

	for _, tt := range testcases {

		rr := httptest.NewRecorder()

		endpoint := fmt.Sprintf("/api/v1/messages/topics/%s?limit=%d&offset=%d", tt.topic, tt.limit, tt.offset)
		request, err := http.NewRequest(http.MethodGet, endpoint, nil)
		assert.NoError(err)

		suite.handler.ServeHTTP(rr, request)

		assert.Equal(tt.expected, rr.Code)
	}
}

func TestMessageControllerSuite(t *testing.T) {
	suite.Run(t, new(MessageControllerSuite))
}
