package spark

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	Listen() error
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

func (e *eventListener) Listen() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("connecting to %s", e.socketURL)

	c, _, err := websocket.DefaultDialer.Dial(e.socketURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	authMsg := []byte(fmt.Sprintf(`{
		"id": %q,
		"type": "authorization",
		"data": {
			"token": "Bearer %s"
		}
	}`, uuid.NewV4(), e.token))

	log.Printf("%s", authMsg)

	if err = c.WriteMessage(websocket.TextMessage, authMsg); err != nil {
		return err
	}

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return c.Close()
		}
	}
}
