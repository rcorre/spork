package main

import (
	"testing"
	"time"

	"github.com/rcorre/spork/mocks"
	"github.com/rcorre/spork/spark"
	"github.com/stretchr/testify/suite"
)

type RoomTestSuite struct {
	suite.Suite
}

func TestRoomTestSuite(t *testing.T) {
	suite.Run(t, new(RoomTestSuite))
}

func (suite *RoomTestSuite) TestTitle() {
	r := NewRoom(&spark.Room{Title: "My Awesome Room"}, nil, nil)
	suite.Equal(r.Title(), "My Awesome Room")
}

func t(s string) *time.Time {
	ret, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		panic(err)
	}
	return &ret
}

func (suite *RoomTestSuite) TestLoad() {
	messageService := &mocks.MessageService{}
	personCache := &mocks.PersonCache{}
	messageService.On(
		"List",
		"abc-123",
	).Return(
		[]spark.Message{
			{PersonID: "ID1", Text: "biz", Created: t("2017-02-01T01:01:01.000Z")},
			{PersonID: "ID1", Text: "foo", Created: t("2017-02-02T02:01:01.000Z")},
			{PersonID: "ID2", Text: "bar", Created: t("2016-01-01T01:01:01.000Z")},
			{PersonID: "ID3", Text: "buz", Created: t("2017-02-02T01:01:01.000Z")},
			{PersonID: "ID2", Text: "baz", Created: t("2017-01-01T01:01:01.000Z")},
		},
		nil,
	)

	personCache.On("Get", "ID1").Return("person1", nil).Times(2)
	personCache.On("Get", "ID2").Return("person2", nil).Times(2)
	personCache.On("Get", "ID3").Return("person3", nil).Times(1)

	room := NewRoom(
		&spark.Room{ID: "abc-123"},
		messageService,
		personCache,
	)
	err := room.Load()
	suite.Nil(err)

	expected := []Message{
		{Sender: "person2", Text: "bar", Time: *t("2016-01-01T01:01:01.000Z")},
		{Sender: "person2", Text: "baz", Time: *t("2017-01-01T01:01:01.000Z")},
		{Sender: "person1", Text: "biz", Time: *t("2017-02-01T01:01:01.000Z")},
		{Sender: "person3", Text: "buz", Time: *t("2017-02-02T01:01:01.000Z")},
		{Sender: "person1", Text: "foo", Time: *t("2017-02-02T02:01:01.000Z")},
	}

	actual := room.Messages()
	suite.Nil(err)
	suite.Equal(expected, actual)

	messageService.AssertExpectations(suite.T())
	personCache.AssertExpectations(suite.T())
}

func (suite *RoomTestSuite) TestSend() {
	messageService := &mocks.MessageService{}
	personCache := &mocks.PersonCache{}
	messageService.On(
		"Post",
		spark.Message{
			RoomID: "abc-123",
			Text:   "tally-ho!",
		},
	).Return(
		spark.Message{
			Text:     "tally-ho!",
			PersonID: "person-123",
			Created:  t("2016-01-01T01:01:01.000Z"),
		},
		nil,
	).Once()

	personCache.On("Get", "person-123").Return("person1", nil).Once()

	room := NewRoom(
		&spark.Room{ID: "abc-123"},
		messageService,
		personCache,
	)
	err := room.Send("tally-ho!")
	suite.Nil(err)

	expected := []Message{{
		Sender: "person1",
		Text:   "tally-ho!",
		Time:   *t("2016-01-01T01:01:01.000Z"),
	}}

	actual := room.Messages()
	suite.Nil(err)
	suite.Equal(expected, actual)

	messageService.AssertExpectations(suite.T())
	personCache.AssertExpectations(suite.T())
}
