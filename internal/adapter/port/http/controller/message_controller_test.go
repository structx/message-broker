package controller_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/trevatk/block-broker/internal/adapter/logging"
	"github.com/trevatk/block-broker/internal/adapter/port/http/controller"
	"github.com/trevatk/block-broker/internal/core/domain"
)

type MessageControllerSuite struct {
	suite.Suite
	mux *http.ServeMux
}

func (suite *MessageControllerSuite) SetupTest() {

	assert := suite.Assert()

	logger, err := logging.NewLogger()
	assert.NoError(err)

	mockMessenger := domain.NewMockMessenger(suite.T())
	mockMessenger.On("Read", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(&domain.Message{}, nil).Once()

	s := fuego.NewServer()

	c := controller.NewMessageController(logger, mockMessenger)
	c.RegisterRoutesV1(s)

	suite.mux = s.Mux
}

func (suite *MessageControllerSuite) TestGet() {

	assert := suite.Assert()

	testcases := []struct {
		uid      string
		expected int
	}{
		{
			uid:      "322befbd-ed13-4566-93e7-24fe87e5306f",
			expected: http.StatusAccepted,
		},
		{
			uid:      "12",
			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range testcases {

		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/message/%s", tt.uid), nil)
		assert.NoError(err)

		suite.mux.ServeHTTP(rr, request)

		assert.Equal(tt.expected, rr.Code)
	}
}

func TestMessageControllerSuite(t *testing.T) {
	suite.Run(t, new(MessageControllerSuite))
}
