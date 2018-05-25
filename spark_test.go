package main

import (
	"encoding/json"
	"io/ioutil"
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
		suite.Equal("GET", r.Method)
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
	expected := []*Room{
		{Title: "Foo"},
		{Title: "Bar"},
		{Title: "Baz"},
	}
	suite.Equal(expected, rooms)
}

func (suite *SparkTestSuite) TestPeople() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal("GET", r.Method)
		suite.Equal("/people", r.URL.Path)
		suite.Equal("Bearer fake-token", r.Header.Get("Authorization"))
		suite.Equal(r.URL.Query().Get("id"), "one,two,three")
		w.Write([]byte(`{
			"items": [
				{"displayName": "Foo"},
				{"displayName": "Bar"},
				{"displayName": "Baz"}
			]
		}`))
	}))
	spark := NewSpark(srv.URL, "fake-token")

	actual, err := spark.People([]string{"one", "two", "three"})
	suite.Nil(err)
	suite.Equal([]*Person{
		{DisplayName: "Foo"},
		{DisplayName: "Bar"},
		{DisplayName: "Baz"},
	}, actual)
}

func (suite *SparkTestSuite) TestMe() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal("GET", r.Method)
		suite.Equal("/people/me", r.URL.Path)
		suite.Equal("Bearer fake-token", r.Header.Get("Authorization"))
		w.Write([]byte(`{
			"id": "mee-123",
			"displayName": "Ryan"
		}`))
	}))
	spark := NewSpark(srv.URL, "fake-token")

	actual, err := spark.Me()
	suite.Nil(err)
	suite.Equal(&Person{
		ID:          "mee-123",
		DisplayName: "Ryan",
	}, actual)
}

func (suite *SparkTestSuite) TestMessages() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal("GET", r.Method)
		suite.Equal("/messages", r.URL.Path)
		suite.Equal("room-12345", r.URL.Query().Get("roomId"))
		suite.Equal("Bearer fake-token", r.Header.Get("Authorization"))
		w.Write([]byte(`{
			"items": [
				{"text": "Foo"},
				{"text": "Bar"},
				{"text": "Baz"}
			]
		}`))
	}))
	spark := NewSpark(srv.URL, "fake-token")

	actual, err := spark.Messages("room-12345")
	suite.Nil(err)
	suite.Equal([]*Message{
		{Text: "Foo"},
		{Text: "Bar"},
		{Text: "Baz"},
	}, actual)
}

func (suite *SparkTestSuite) TestSend() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Equal("POST", r.Method)
		suite.Equal("/messages", r.URL.Path)
		suite.Equal("Bearer fake-token", r.Header.Get("Authorization"))

		defer r.Body.Close()
		b, _ := ioutil.ReadAll(r.Body)
		var msg Message
		suite.NoError(json.Unmarshal(b, &msg))
		suite.Equal(Message{RoomID: "abc-123", Text: "foobar"}, msg)

		w.Write([]byte(`{
			"id": "xyz-345",
			"roomID": "abc-123",
			"text": "foobar"
		}`))
	}))
	spark := NewSpark(srv.URL, "fake-token")

	actual, err := spark.Send(&Message{
		RoomID: "abc-123",
		Text:   "foobar",
	})
	suite.Nil(err)
	suite.Equal(&Message{
		ID:     "xyz-345",
		RoomID: "abc-123",
		Text:   "foobar",
	}, actual)
}
