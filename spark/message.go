package spark

type MessageService interface {
	List(roomID string) ([]Message, error)
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
	Created         string
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
