package spark

import "time"

type MessageService interface {
	List(roomID string) ([]Message, error)
	Post(msg Message) (Message, error)
}

type Message struct {
	ID              string
	RoomID          string
	RoomType        string
	ToPersonID      string
	ToPersonEmail   string
	Text            string
	Markdown        string
	Files           []string
	PersonID        string
	PersonEmail     string
	Created         time.Time
	MentionedPeople []string
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

// Post posts a message
// roomID is the id of the room to list messages from
func (svc *messageService) Post(msg Message) (Message, error) {
	var out Message
	err := svc.rest.Post("messages", &msg, &out)
	return out, err
}
