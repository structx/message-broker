package application_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/trevatk/block-broker/internal/core/application"
	"github.com/trevatk/block-broker/internal/core/domain"
)

type MessagingServiceSuite struct {
	suite.Suite
	m domain.Messenger
}

func (suite *MessagingServiceSuite) SetupTest() {

	mockChain := domain.NewMockChain(suite.T())
	mockChain.On(
		"AddTx",
		mock.AnythingOfType("string"), // data
		mock.AnythingOfType("string"), // action
		mock.Anything,                 // payload
		mock.AnythingOfType("string"), // signature
	).
		Return("0cf6b052518a08ed5299676530c3d57f2c94dbe9ec26184b4b9d3baf", nil).
		Maybe()

	mockChain.On(
		"ReadTx",
		mock.Anything, // hash
	).
		Return(&domain.Tx{
			ID:        []byte(""),
			Topic:     "unit.test",
			Action:    "publish",
			Payload:   []byte("hello world"),
			Timestamp: time.Now().Format(time.RFC3339),
			Sig:       "xoxo",
		}, nil).
		Maybe()

	mockChain.On(
		"ListTransactions",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
	).
		Return([]*domain.Tx{}, nil).
		Maybe()

	mockChain.On(
		"ListTransactionsByAction",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
	).
		Return([]*domain.Tx{}, nil).
		Maybe()

	suite.m = application.NewMessagingService(mockChain)
}

func (suite *MessagingServiceSuite) TestCreate() {

	assert := suite.Assert()

	testcases := []struct {
		expected   error
		newMessage *domain.NewMessage
	}{
		{
			expected: nil,
			newMessage: &domain.NewMessage{
				Topic:     "unit.test",
				Payload:   []byte("hello world"),
				Signature: "golang",
			},
		},
	}

	for _, tt := range testcases {

		msg, err := suite.m.Create(tt.newMessage)
		assert.Equal(tt.expected, err)

		assert.Equal(tt.newMessage.Topic, msg.Topic)
		assert.Equal(tt.newMessage.Payload, msg.Payload)
		assert.NotNil(msg.CreatedAt)
	}
}

func (suite *MessagingServiceSuite) TestRead() {

	assert := suite.Assert()

	testcases := []struct {
		hash     string
		expected error
	}{
		{
			hash:     "0cf6b052518a08ed5299676530c3d57f2c94dbe9ec26184b4b9d3baf",
			expected: nil,
		},
	}

	for _, tt := range testcases {

		_, err := suite.m.Read(tt.hash)
		assert.Equal(tt.expected, err)
	}
}

func (suite *MessagingServiceSuite) TestList() {

	assert := suite.Assert()

	testcases := []struct {
		limit, offset int
		expected      error
	}{
		{
			expected: nil,
		},
	}

	for _, tt := range testcases {

		_, err := suite.m.List(tt.limit, tt.offset)
		assert.Equal(tt.expected, err)
	}
}

func (suite *MessagingServiceSuite) TestListByTopic() {

	assert := suite.Assert()

	testcases := []struct {
		expected      error
		topic         string
		limit, offset int
	}{
		{
			expected: nil,
			topic:    "unit.test",
			limit:    10,
			offset:   0,
		},
	}

	for _, tt := range testcases {

		_, err := suite.m.ListByTopic(tt.topic, tt.limit, tt.offset)
		assert.Equal(tt.expected, err)
	}
}

func (suite *MessagingServiceSuite) TestListTopics() {

	assert := suite.Assert()

	testcases := []struct {
		limit, offset int
		expected      error
	}{
		{
			expected: nil,
		},
	}

	for _, tt := range testcases {

		_, err := suite.m.ListTopics(tt.limit, tt.offset)
		assert.Equal(tt.expected, err)
	}
}

func TestMessagingServiceSuite(t *testing.T) {
	suite.Run(t, new(MessagingServiceSuite))
}
