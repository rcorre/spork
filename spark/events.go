package spark

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/romana/rlog"
	uuid "github.com/satori/go.uuid"
)

const socketURL = "https://wdm-a.wbx2.com/wdm/api/v1/devices"

type dialer interface {
	Dial(urlStr string, requestHeader http.Header) (connection, *http.Response, error)
}

type connection interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

// Cisco Spark has a websocket interface to listen for message events
// It isn't documented, I found it here:
// https://github.com/marchfederico/ciscospark-websocket-events
type EventListener interface {
	Devices() (interface{}, error)
	Register() error
	UnRegister() error
	Listen() (chan Event, chan error, error)
}

type eventListener struct {
	rest      RESTClient
	token     string
	deviceURL string
	socketURL string
	conn      connection
	connect   func(url string) (connection, error)
}

func NewEventListener(token string) EventListener {
	return &eventListener{
		token: token,
		rest:  NewRESTClient(socketURL, token),
		connect: func(url string) (connection, error) {
			rlog.Debugf("connecting to %s", url)

			c, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				return nil, err
			}
			return connection(c), nil
		},
	}
}

func (e *eventListener) Devices() (interface{}, error) {
	var out interface{}
	err := e.rest.Get("", nil, &out)
	return out, err
}

func (e *eventListener) Register() error {
	var resp struct {
		URL          string
		WebSocketURL string
	}

	spec := map[string]string{
		"deviceName":     "spork",
		"deviceType":     "DESKTOP",
		"localizedModel": "go",
		"model":          "go",
		"name":           "spork",
		"systemName":     "spork",
		"systemVersion":  "0.1",
	}

	err := e.rest.Post("", spec, &resp)
	e.deviceURL = resp.URL
	e.socketURL = resp.WebSocketURL
	return err
}

func (e *eventListener) UnRegister() error {
	if e.deviceURL == "" {
		return nil
	}

	req, err := http.NewRequest("DELETE", e.deviceURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+e.token)
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Failed to unregister device", resp)
	}
	return err
}

type Event struct {
	ID   string
	Data struct {
		EventType string
		Activity  struct {
			Verb  string
			Actor struct {
				ID           string
				ObjectType   string
				DisplayName  string
				OrgId        string
				EmailAddress string
				EntryUUID    string
				Type         string
			}
		}
	}
}

func (e *eventListener) Listen() (chan Event, chan error, error) {
	rlog.Infof("connecting to %s", e.socketURL)

	if e.socketURL == "" {
		return nil, nil, fmt.Errorf("Cannot Listen() before Register()")
	}

	c, err := e.connect(e.socketURL)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to dial websocket: %v", err)
	}

	evChan := make(chan Event)
	errChan := make(chan error)

	go func() {
		defer c.Close()
		defer close(evChan)
		defer close(errChan)

		rlog.Debug("Seding auth to websocket...")

		authMsg := []byte(fmt.Sprintf(`{
			"id": %q,
			"type": "authorization",
			"data": {
				"token": "Bearer %s"
			}
		}`, uuid.NewV4(), e.token))

		if err = c.WriteMessage(websocket.TextMessage, authMsg); err != nil {
			errChan <- err
			return
		}

		rlog.Debug("Websocket auth successful!")

		for {
			_, message, err := c.ReadMessage()
			rlog.Debugf("Websocket recv msg: %s, err: %v", message, err)

			if err != nil {
				errChan <- err
				return
			}
			var ev Event
			if err := json.Unmarshal(message, &ev); err != nil {
				errChan <- err
			} else {
				evChan <- ev
			}
		}
	}()

	return evChan, errChan, nil
}
