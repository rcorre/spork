package spark

const defaultURL = "https://api.ciscospark.com/v1/"

type Client struct {
	Rooms    RoomService
	Messages MessageService
	People   PeopleService

	Events EventListener
}

func New(url, token string) *Client {
	if url == "" {
		url = defaultURL
	}

	rest := NewRESTClient(url, token)

	return &Client{
		Rooms:    &roomService{rest: rest},
		Messages: &messageService{rest: rest},
		People:   &peopleService{rest: rest},
		Events:   NewEventListener(token),
	}
}
