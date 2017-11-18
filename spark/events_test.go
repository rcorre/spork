package spark

import (
	"encoding/json"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const exampleEvent = `{
	"id": "4",
	"data": {
		"activity": {
			"verb": "acknowledge",
			"actor": {
				"id": "me@example.com",
				"objectType": "person",
				"displayName": "Myself",
				"orgId": "a6b96c03-0aed-495e-9f5b-6c92ab21c9e8",
				"emailAddress": "me@example.com",
				"entryUUID": "a6b96c03-0aed-495e-9f5b-6c92ab21c9e8",
				"type": "PERSON"
			}
		},
		"eventType": "conversation.activity"
	}
}`

type EventsTestSuite struct {
	suite.Suite
}

func TestEventsTestSuite(t *testing.T) {
	suite.Run(t, new(EventsTestSuite))
}

type MockConnection struct {
	mock.Mock
}

func (m *MockConnection) ReadMessage() (messageType int, p []byte, err error) {
	args := m.Called()
	return args.Int(0), args.Get(1).([]byte), args.Error(2)
}

func (m *MockConnection) WriteMessage(messageType int, data []byte) error {
	args := m.Called(messageType, data)
	return args.Error(0)
}

func (m *MockConnection) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (suite *EventsTestSuite) TestRegister() {
	restClient := &RESTClientMock{}
	restClient.On(
		"Post",
		"",
		map[string]string{
			"deviceName":     "spork",
			"deviceType":     "DESKTOP",
			"localizedModel": "go",
			"model":          "go",
			"name":           "spork",
			"systemName":     "spork",
			"systemVersion":  "0.1",
		},
		mock.Anything,
	).Run(func(args mock.Arguments) {
		out := args.Get(2).(*struct {
			URL          string
			WebSocketURL string
		})
		out.URL = "http://example.com/device"
		out.WebSocketURL = "http://example.com/socket"
	}).Return(nil)
	listener := eventListener{
		rest: restClient,
	}

	err := listener.Register()
	suite.Nil(err)
	suite.Equal("http://example.com/device", listener.deviceURL)
	suite.Equal("http://example.com/socket", listener.socketURL)
	restClient.AssertExpectations(suite.T())
}

func (suite *EventsTestSuite) TestUnRegister() {
	// TODO
}

func (suite *EventsTestSuite) TestListen() {
	conn := MockConnection{}
	conn.On(
		"WriteMessage",
		websocket.TextMessage,
		mock.MatchedBy(func(data []byte) bool {
			var out struct {
				ID   string
				Type string
				Data map[string]string
			}
			err := json.Unmarshal(data, &out)
			suite.Nil(err)
			return true
		}),
	).Return(nil)

	conn.On("ReadMessage").Return(1, []byte(exampleEvent), nil)
	conn.On("Close").Return(nil)

	e := eventListener{
		socketURL: "wss:example.com",
		token:     "mytoken",
		connect: func(url string) (connection, error) {
			suite.Equal(url, "wss:example.com")
			return &conn, nil
		},
	}
	evChan, errChan, err := e.Listen()
	suite.Nil(err)
	suite.NotNil(evChan)
	suite.NotNil(errChan)
	ev := <-evChan

	suite.Equal(ev.ID, "4")
	suite.Equal(ev.Data.Activity.Verb, "acknowledge")
	suite.Equal(ev.Data.Activity.Actor.ID, "me@example.com")
	suite.Equal(ev.Data.Activity.Actor.ObjectType, "person")
	suite.Equal(ev.Data.Activity.Actor.DisplayName, "Myself")
	suite.Equal(ev.Data.Activity.Actor.OrgId, "a6b96c03-0aed-495e-9f5b-6c92ab21c9e8")
	suite.Equal(ev.Data.Activity.Actor.EmailAddress, "me@example.com")
	suite.Equal(ev.Data.Activity.Actor.EntryUUID, "a6b96c03-0aed-495e-9f5b-6c92ab21c9e8")
	suite.Equal(ev.Data.Activity.Actor.Type, "PERSON")
	suite.Equal(ev.Data.EventType, "conversation.activity")
	suite.NotNil(errChan)
}
