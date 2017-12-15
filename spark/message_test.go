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
			{Text: "Foo"},
			{Text: "Bar"},
			{Text: "Baz"},
		}
	}).Return(nil)
	roomService := NewMessageService(restClient)

	rooms, err := roomService.List("room-12345")
	suite.Nil(err)
	suite.Equal(rooms, []Message{
		{Text: "Foo"},
		{Text: "Bar"},
		{Text: "Baz"},
	})
}

func (suite *MessageTestSuite) TestGet() {
	restClient := &RESTClientMock{}
	restClient.On(
		"Get",
		"messages/msg-12345",
		map[string]string(nil),
		&Message{},
	).Run(func(args mock.Arguments) {
		out := args.Get(2).(*Message)
		*out = Message{Text: "Foo"}
	}).Return(nil)
	svc := NewMessageService(restClient)

	actual, err := svc.Get("msg-12345")
	suite.Nil(err)
	suite.Equal(Message{Text: "Foo"}, actual)
}

func (suite *MessageTestSuite) TestPost() {
	restClient := &RESTClientMock{}
	restClient.On(
		"Post",
		"messages",
		&Message{RoomID: "abc-123", Text: "foobar"},
		&Message{},
	).Run(func(args mock.Arguments) {
		out := args.Get(2).(*Message)
		*out = Message{
			ID:     "xyz-345",
			RoomID: "abc-123",
			Text:   "foobar",
		}
	}).Return(nil)
	roomService := NewMessageService(restClient)

	out, err := roomService.Post(Message{
		RoomID: "abc-123",
		Text:   "foobar",
	})
	suite.Nil(err)
	suite.Equal(out, Message{
		ID:     "xyz-345",
		RoomID: "abc-123",
		Text:   "foobar",
	})
}
