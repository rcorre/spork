package spark

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

const socketURL = "https://wdm-a.wbx2.com/wdm/api/v1/devices"

// Cisco Spark has a websocket interface to listen for message events
// It isn't documented, I found it here:
// https://github.com/marchfederico/ciscospark-websocket-events
type EventListener interface {
	Devices() (interface{}, error)
	Register() error
	UnRegister() error
	Listen() (chan string, chan error, error)
}

type eventListener struct {
	rest      RESTClient
	token     string
	deviceURL string
	socketURL string
}

func NewEventListener(token string) EventListener {
	return &eventListener{
		token: token,
		rest:  NewRESTClient(socketURL, token),
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

func (e *eventListener) Listen() (chan string, chan error, error) {
	log.Printf("connecting to %s", e.socketURL)

	c, _, err := websocket.DefaultDialer.Dial(e.socketURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to dial websocket: %v", err)
	}

	msgChan := make(chan string)
	errChan := make(chan error)

	go func() {
		defer c.Close()
		defer close(msgChan)
		defer close(errChan)

		log.Printf("Websocket auth...")

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

		log.Printf("Websocket auth complete")

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			msgChan <- string(message)
		}
	}()

	return msgChan, errChan, nil
}
