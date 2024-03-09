package application_test

import (
	"testing"

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
	mockChain.EXPECT().AddTx(mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("[]byte"), mock.AnythingOfType("string")).Return("", nil).Once()

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

func TestMessagingServiceSuite(t *testing.T) {
	suite.Run(t, new(MessagingServiceSuite))
}
