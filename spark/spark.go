package spark

const defaultURL = "https://api.ciscospark.com/v1/"
const defaultDeviceURL = "https://wdm-a.wbx2.com/wdm/api/v1/devices"

type Client struct {
	Rooms    RoomService
	Messages MessageService
	People   PeopleService

	Events EventListener
}

// New creates a new Spark client
// url is the spark api url (default "https://api.ciscospark.com/v1/")
// deviceURL is the device registration url (default "https://wdm-a.wbx2.com/wdm/api/v1/devices")
// token is the spark API token
func New(url, deviceURL, token string) *Client {
	if url == "" {
		url = defaultURL
	}

	if deviceURL == "" {
		deviceURL = defaultDeviceURL
	}

	rest := NewRESTClient(url, token)

	return &Client{
		Rooms:    &roomService{rest: rest},
		Messages: &messageService{rest: rest},
		People:   &peopleService{rest: rest},
		Events:   NewEventListener(deviceURL, token),
	}
}
