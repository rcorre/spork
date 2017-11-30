package spark

import "time"

type RoomService interface {
	List() ([]Room, error)
}

type Room struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Type         string    `json:"type"`
	IsLocked     bool      `json:"isLocked"`
	TeamID       string    `json:"teamId"`
	LastActivity time.Time `json:"lastActivity"`
	Created      time.Time `json:"created"`
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
