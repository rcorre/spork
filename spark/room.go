package spark

import "time"

type RoomService interface {
	List() ([]Room, error)
}

type Room struct {
	ID           string
	Title        string
	Type         string
	IsLocked     bool
	TeamId       string
	LastActivity time.Time
	Created      time.Time
}

type roomService struct {
	rest RESTClient
}

func NewRoomService(rest RESTClient) RoomService {
	return &roomService{rest: rest}
}

func (svc *roomService) List() ([]Room, error) {
	var list struct {
		Items []Room
	}
	err := svc.rest.Get("rooms", map[string]string{}, &list)
	return list.Items, err
}
