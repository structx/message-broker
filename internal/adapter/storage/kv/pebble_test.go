package kv_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trevatk/block-broker/internal/adapter/setup"
	"github.com/trevatk/block-broker/internal/adapter/storage/kv"
)

func init() {
	_ = os.Setenv("KV_DIR", "./testfiles")
}

type PebbleSuite struct {
	suite.Suite
	pebble *kv.Pebble
}

func (suite *PebbleSuite) SetupTest() {

	assert := suite.Assert()
	ctx := context.TODO()

	cfg := setup.NewConfig()
	assert.NoError(setup.ProcessConfigWithEnv(ctx, cfg))

	var err error

	suite.pebble, err = kv.NewPebble(cfg)
	assert.NoError(err)
}

func (suite *PebbleSuite) TestPut() {

	assert := suite.Assert()

	err := suite.pebble.Put([]byte("0"), []byte("hello world"))
	assert.NoError(err)
}

func TestPebbleSuite(t *testing.T) {
	suite.Run(t, new(PebbleSuite))
}
