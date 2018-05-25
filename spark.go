package main

import (
	"strings"
	"time"
)

// Spark is an interface to the cisco spark API
// see https://developer.ciscospark.com/getting-started.html
type Spark interface {
	Rooms() ([]*Room, error)
	People(ids []string) ([]*Person, error)
	Me() (*Person, error)
	Messages(roomId string) ([]*Message, error)
	Send(msg *Message) (*Message, error)
}

// Room is a spark room
type Room struct {
	ID           string     `json:"id,omitempty"`
	Title        string     `json:"title,omitempty"`
	Type         string     `json:"type,omitempty"`
	IsLocked     bool       `json:"isLocked,omitempty"`
	TeamID       string     `json:"teamId,omitempty"`
	LastActivity *time.Time `json:"lastActivity,omitempty"`
	Created      *time.Time `json:"created,omitempty"`
}

// Person is a spark user
type Person struct {
	ID            string     `json:"id,omitempty"`
	Emails        []string   `json:"emails,omitempty"`
	DisplayName   string     `json:"displayName,omitempty"`
	FirstName     string     `json:"firstName,omitempty"`
	LastName      string     `json:"lastName,omitempty"`
	Avatar        string     `json:"avatar,omitempty"`
	OrgID         string     `json:"orgId,omitempty"`
	Roles         []string   `json:"roles,omitempty"`
	Licenses      []string   `json:"licenses,omitempty"`
	Created       *time.Time `json:"created,omitempty"`
	Timezone      string     `json:"timezone,omitempty"`
	LastActivity  *time.Time `json:"lastActivity,omitempty"`
	Status        string     `json:"status,omitempty"`
	InvitePending bool       `json:"invitePending,omitempty"`
	LoginEnabled  bool       `json:"loginEnabled,omitempty"`
}

// Message is a spark message
type Message struct {
	ID              string     `json:"id,omitempty"`
	RoomID          string     `json:"roomId,omitempty"`
	RoomType        string     `json:"roomType,omitempty"`
	ToPersonID      string     `json:"toPersonId,omitempty"`
	ToPersonEmail   string     `json:"toPersonEmail,omitempty"`
	Text            string     `json:"text,omitempty"`
	Markdown        string     `json:"markdown,omitempty"`
	HTML            string     `json:"html,omitempty"`
	Files           []string   `json:"files,omitempty"`
	PersonID        string     `json:"personId,omitempty"`
	PersonEmail     string     `json:"personEmail,omitempty"`
	Created         *time.Time `json:"created,omitempty"`
	MentionedPeople []string   `json:"mentionedPeople,omitempty"`
}

type spark struct {
	url   string
	token string
	rest  RESTClient
}

// New creates a new Spark client
// url is the spark api url
// token is the spark API token
func NewSpark(url, token string) Spark {
	return &spark{
		url:   url,
		token: token,
		rest:  NewRESTClient(url, token),
	}
}

func (s *spark) Rooms() ([]*Room, error) {
	var list struct {
		Items []*Room
	}
	err := s.rest.Get("rooms", map[string]string{}, &list)
	return list.Items, err
}

// People looks up people by ID
// ids are the IDs of people to list
func (s *spark) People(ids []string) ([]*Person, error) {
	var list struct {
		Items []*Person
	}
	params := map[string]string{
		"id": strings.Join(ids, ","),
	}
	err := s.rest.Get("people", params, &list)
	return list.Items, err
}

// Me returns the Person representing the current user
func (s *spark) Me() (*Person, error) {
	var me Person
	err := s.rest.Get("people/me", nil, &me)
	return &me, err
}

// List lists messages from a room
// roomID is the id of the room to list messages from
func (s *spark) Messages(roomID string) ([]*Message, error) {
	var list struct {
		Items []*Message
	}
	params := map[string]string{
		"roomId": roomID,
	}
	err := s.rest.Get("messages", params, &list)
	return list.Items, err
}

// Send sends a message
// msg is the message to post
func (s *spark) Send(msg *Message) (*Message, error) {
	var out Message
	err := s.rest.Post("messages", &msg, &out)
	return &out, err
}
