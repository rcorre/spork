package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SparkTestSuite struct {
	suite.Suite
}

func TestSparkTestSuite(t *testing.T) {
	suite.Run(t, new(SparkTestSuite))
}

func (suite *SparkTestSuite) TestRooms() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal("/rooms", r.URL.Path)
		suite.Equal("Bearer fake-token", r.Header.Get("Authorization"))
		w.Write([]byte(`{
			"items": [
				{"title": "Foo"},
				{"title": "Bar"},
				{"title": "Baz"}
			]
		}`))
	}))
	spark := NewSpark(srv.URL, "fake-token")

	rooms, err := spark.Rooms()
	suite.Nil(err)
	suite.Equal(rooms, []*Room{
		{Title: "Foo"},
		{Title: "Bar"},
		{Title: "Baz"},
	})
}

//func (suite *SparkTestSuite) TestPeople() {
//	restClient := &RESTClientMock{}
//	restClient.On(
//		"Get",
//		"people",
//		map[string]string{
//			"id": "one,two,three",
//		},
//		&struct{ Items []Person }{},
//	).Run(func(args mock.Arguments) {
//		out := args.Get(2).(*struct{ Items []Person })
//		out.Items = []Person{
//			{DisplayName: "Foo"},
//			{DisplayName: "Bar"},
//			{DisplayName: "Baz"},
//		}
//	}).Return(nil)
//	peopleService := NewPeopleService(restClient)
//
//	rooms, err := peopleService.List([]string{"one", "two", "three"})
//	suite.Nil(err)
//	suite.Equal(rooms, []Person{
//		{DisplayName: "Foo"},
//		{DisplayName: "Bar"},
//		{DisplayName: "Baz"},
//	})
//
//	restClient.AssertExpectations(suite.T())
//}
//
//func (suite *SparkTestSuite) TestMe() {
//	restClient := &RESTClientMock{}
//	restClient.On(
//		"Get",
//		"people/me",
//		map[string]string(nil),
//		mock.Anything,
//	).Run(func(args mock.Arguments) {
//		out := args.Get(2).(*Person)
//		*out = Person{
//			ID:          "mee-123",
//			DisplayName: "Ryan",
//		}
//	}).Return(nil)
//	peopleService := NewPeopleService(restClient)
//
//	me, err := peopleService.Me()
//	suite.Nil(err)
//	suite.Equal(me.ID, "mee-123")
//	suite.Equal(me.DisplayName, "Ryan")
//
//	restClient.AssertExpectations(suite.T())
//}
//
//func (suite *SparkTestSuite) TestMessages() {
//	restClient := &RESTClientMock{}
//	restClient.On(
//		"Get",
//		"messages",
//		map[string]string{"roomId": "room-12345"},
//		&struct{ Items []Message }{},
//	).Run(func(args mock.Arguments) {
//		out := args.Get(2).(*struct{ Items []Message })
//		out.Items = []Message{
//			{Text: "Foo"},
//			{Text: "Bar"},
//			{Text: "Baz"},
//		}
//	}).Return(nil)
//	s := NewSpark("", "")
//
//	rooms, err := roomService.List("room-12345")
//	suite.Nil(err)
//	suite.Equal(rooms, []Message{
//		{Text: "Foo"},
//		{Text: "Bar"},
//		{Text: "Baz"},
//	})
//}
//
//func (suite *SparkTestSuite) TestSend() {
//	restClient := &RESTClientMock{}
//	restClient.On(
//		"Post",
//		"messages",
//		&Message{RoomID: "abc-123", Text: "foobar"},
//		&Message{},
//	).Run(func(args mock.Arguments) {
//		out := args.Get(2).(*Message)
//		*out = Message{
//			ID:     "xyz-345",
//			RoomID: "abc-123",
//			Text:   "foobar",
//		}
//	}).Return(nil)
//	roomService := NewMessageService(restClient)
//
//	out, err := roomService.Post(Message{
//		RoomID: "abc-123",
//		Text:   "foobar",
//	})
//	suite.Nil(err)
//	suite.Equal(out, Message{
//		ID:     "xyz-345",
//		RoomID: "abc-123",
//		Text:   "foobar",
//	})
//}
