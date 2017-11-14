package spark

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MessageTestSuite struct {
	suite.Suite
}

func TestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}

func (suite *MessageTestSuite) TestList() {
	restClient := &RESTClientMock{}
	restClient.On(
		"Get",
		"messages",
		map[string]string{"roomId": "room-12345"},
		&struct{ Items []Message }{},
	).Run(func(args mock.Arguments) {
		out := args.Get(2).(*struct{ Items []Message })
		out.Items = []Message{
			Message{Text: "Foo"},
			Message{Text: "Bar"},
			Message{Text: "Baz"},
		}
	}).Return(nil)
	roomService := NewMessageService(restClient)

	rooms, err := roomService.List("room-12345")
	suite.Nil(err)
	suite.Equal(rooms, []Message{
		Message{Text: "Foo"},
		Message{Text: "Bar"},
		Message{Text: "Baz"},
	})
}
