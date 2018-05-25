package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SparkMock struct {
	mock.Mock
}

func (m *SparkMock) Get(path string, params map[string]string, out interface{}) error {
	args := m.Called(path, params, out)
	return args.Error(0)
}

func (m *SparkMock) Rooms() ([]*Room, error) {
	args := m.Called()
	return args.Get(0).([]*Room), args.Error(1)
}

func (m *SparkMock) People(ids []string) ([]*Person, error) {
	args := m.Called(ids)
	return args.Get(0).([]*Person), args.Error(1)
}

func (m *SparkMock) Me() (*Person, error) {
	args := m.Called()
	return args.Get(0).(*Person), args.Error(1)
}

func (m *SparkMock) Messages(roomID string) ([]*Message, error) {
	args := m.Called(roomID)
	return args.Get(0).([]*Message), args.Error(1)
}

func (m *SparkMock) Send(msg *Message) (*Message, error) {
	args := m.Called(msg)
	return args.Get(0).(*Message), args.Error(1)
}

type CoreTestSuite struct {
	suite.Suite
}

func TestCoreTestSuite(t *testing.T) {
	suite.Run(t, new(CoreTestSuite))
}

func t(s string) *time.Time {
	ret, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		panic(err)
	}
	return &ret
}

func (suite *CoreTestSuite) TestLoadRooms() {
	spark := SparkMock{}
	spark.On("Rooms").Return([]*Room{
		{Title: "foo", LastActivity: &time.Time{}},
		{Title: "bar", LastActivity: &time.Time{}},
		{Title: "baz", LastActivity: &time.Time{}},
	}, nil)
	c := Core{
		spark: &spark,
	}

	err := c.LoadRooms()
	suite.Nil(err)
}

//func (suite *CoreTestSuite) TestLoadMessages() {
//	messageService := &mocks.MessageService{}
//	personCache := &mocks.PersonCache{}
//	messageService.On(
//		"List",
//		"abc-123",
//	).Return(
//		[]spark.Message{
//			{PersonID: "ID1", Text: "biz", Created: t("2017-02-01T01:01:01.000Z")},
//			{PersonID: "ID1", Text: "foo", Created: t("2017-02-02T02:01:01.000Z")},
//			{PersonID: "ID2", Text: "bar", Created: t("2016-01-01T01:01:01.000Z")},
//			{PersonID: "ID3", Text: "buz", Created: t("2017-02-02T01:01:01.000Z")},
//			{PersonID: "ID2", Text: "baz", Created: t("2017-01-01T01:01:01.000Z")},
//		},
//		nil,
//	)
//
//	personCache.On("Get", "ID1").Return("person1", nil).Times(2)
//	personCache.On("Get", "ID2").Return("person2", nil).Times(2)
//	personCache.On("Get", "ID3").Return("person3", nil).Times(1)
//
//	room := NewCore(
//		&spark.Core{ID: "abc-123"},
//		messageService,
//		personCache,
//	)
//	err := room.Load()
//	suite.Nil(err)
//
//	expected := []Message{
//		{Sender: "person2", Text: "bar", Time: *t("2016-01-01T01:01:01.000Z")},
//		{Sender: "person2", Text: "baz", Time: *t("2017-01-01T01:01:01.000Z")},
//		{Sender: "person1", Text: "biz", Time: *t("2017-02-01T01:01:01.000Z")},
//		{Sender: "person3", Text: "buz", Time: *t("2017-02-02T01:01:01.000Z")},
//		{Sender: "person1", Text: "foo", Time: *t("2017-02-02T02:01:01.000Z")},
//	}
//
//	actual := room.Messages()
//	suite.Nil(err)
//	suite.Equal(expected, actual)
//
//	messageService.AssertExpectations(suite.T())
//	personCache.AssertExpectations(suite.T())
//}
//
//func (suite *CoreTestSuite) TestSend() {
//	messageService := &mocks.MessageService{}
//	personCache := &mocks.PersonCache{}
//	messageService.On(
//		"Post",
//		spark.Message{
//			CoreID: "abc-123",
//			Text:   "tally-ho!",
//		},
//	).Return(
//		spark.Message{
//			Text:     "tally-ho!",
//			PersonID: "person-123",
//			Created:  t("2016-01-01T01:01:01.000Z"),
//		},
//		nil,
//	).Once()
//
//	personCache.On("Get", "person-123").Return("person1", nil).Once()
//
//	room := NewCore(
//		&spark.Core{ID: "abc-123"},
//		messageService,
//		personCache,
//	)
//	err := room.Send("tally-ho!")
//	suite.Nil(err)
//
//	expected := []Message{{
//		Sender: "person1",
//		Text:   "tally-ho!",
//		Time:   *t("2016-01-01T01:01:01.000Z"),
//	}}
//
//	actual := room.Messages()
//	suite.Nil(err)
//	suite.Equal(expected, actual)
//
//	messageService.AssertExpectations(suite.T())
//	personCache.AssertExpectations(suite.T())
//}
