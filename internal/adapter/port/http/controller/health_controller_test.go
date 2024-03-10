package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trevatk/block-broker/internal/adapter/logging"
	"github.com/trevatk/block-broker/internal/adapter/port/http/router"
)

type HealthSuite struct {
	suite.Suite
	handler http.Handler
}

func (suite *HealthSuite) SetupTest() {

	assert := suite.Assert()

	logger, err := logging.NewLogger()
	assert.NoError(err)

	suite.handler = router.NewRouter(logger, nil)
}

func (suite *HealthSuite) TestHealthCheck() {

	assert := suite.Assert()

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/health", nil)
	assert.NoError(err)

	suite.handler.ServeHTTP(rr, request)

	assert.Equal(http.StatusOK, rr.Code)
}

func TestHealthSuite(t *testing.T) {
	suite.Run(t, new(HealthSuite))
}
