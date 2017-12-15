package spark

import (
	"fmt"
	"time"
)

type MessageService interface {
	List(roomID string) ([]Message, error)
	Get(id string) (Message, error)
	Post(msg Message) (Message, error)
}

type Message struct {
	ID              string     `json:"id,omitempty"`
	RoomID          string     `json:"roomId,omitempty"`
	RoomType        string     `json:"roomType,omitempty"`
	ToPersonID      string     `json:"toPersonId,omitempty"`
	ToPersonEmail   string     `json:"toPersonEmail,omitempty"`
	Text            string     `json:"text,omitempty"`
	Markdown        string     `json:"markdown,omitempty"`
	Files           []string   `json:"files,omitempty"`
	PersonID        string     `json:"personId,omitempty"`
	PersonEmail     string     `json:"personEmail,omitempty"`
	Created         *time.Time `json:"created,omitempty"`
	MentionedPeople []string   `json:"mentionedPeople,omitempty"`
}

type messageService struct {
	rest RESTClient
}

func NewMessageService(rest RESTClient) MessageService {
	return &messageService{rest: rest}
}

// List lists messages from a room
// roomID is the id of the room to list messages from
func (svc *messageService) List(roomID string) ([]Message, error) {
	var list struct {
		Items []Message
	}
	params := map[string]string{
		"roomId": roomID,
	}
	err := svc.rest.Get("messages", params, &list)
	return list.Items, err
}

// Get retrieves a single message
// id is the message id
func (svc *messageService) Get(id string) (Message, error) {
	var msg Message
	url := fmt.Sprintf("messages/%s", id)
	err := svc.rest.Get(url, nil, &msg)
	return msg, err
}

// Post posts a message
// roomID is the id of the room to list messages from
func (svc *messageService) Post(msg Message) (Message, error) {
	var out Message
	err := svc.rest.Post("messages", &msg, &out)
	return out, err
}
