package spark

import "time"

type RoomService interface {
	List() ([]Room, error)
}

type Room struct {
	ID           string     `json:"id,omitempty"`
	Title        string     `json:"title,omitempty"`
	Type         string     `json:"type,omitempty"`
	IsLocked     bool       `json:"isLocked,omitempty"`
	TeamID       string     `json:"teamId,omitempty"`
	LastActivity *time.Time `json:"lastActivity,omitempty"`
	Created      *time.Time `json:"created,omitempty"`
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
