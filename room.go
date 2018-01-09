package main

import (
	"sort"
	"time"

	"github.com/rcorre/spork/spark"
)

// Room represents a spark room
type Room interface {
	ID() string
	Title() string
	LastActivity() time.Time
	Load() error
	Messages() []Message
	Send(text string) error
}

type room struct {
	data     spark.Room
	svc      spark.MessageService
	people   PersonCache
	messages MessageList
}

// NewRoom creates a Room wrapping a spark.Room
func NewRoom(src *spark.Room, svc spark.MessageService, people PersonCache) Room {
	return &room{
		data:   *src,
		svc:    svc,
		people: people,
	}
}

func (r *room) ID() string {
	return r.data.ID
}

func (r *room) Title() string {
	return r.data.Title
}

func (r *room) LastActivity() time.Time {
	return *r.data.LastActivity
}

func (r *room) Messages() []Message {
	return r.messages
}

func (r *room) Load() error {
	messages, err := r.svc.List(r.data.ID)
	if err != nil {
		return err
	}

	r.messages = make(MessageList, len(messages))
	for i, msg := range messages {
		sender, err := r.people.Get(msg.PersonID)
		if err != nil {
			return err
		}

		r.messages[i] = Message{
			Text:   msg.Text,
			HTML:   msg.HTML,
			Sender: sender,
			Time:   *msg.Created,
		}
	}

	r.messages.Sort()
	return nil
}

func (r *room) Send(text string) error {
	msg, err := r.svc.Post(spark.Message{
		RoomID: r.data.ID,
		Text:   text,
	})
	if err != nil {
		return err
	}

	sender, err := r.people.Get(msg.PersonID)
	if err != nil {
		return err
	}

	r.messages = append(r.messages, Message{
		Text:   msg.Text,
		Sender: sender,
		Time:   *msg.Created,
	})
	return nil
}

// RoomList implements sort.Interface to sort rooms by last activity
type RoomList []Room

func (m RoomList) Len() int {
	return len(m)
}

func (m RoomList) Less(i, j int) bool {
	return m[i].LastActivity().After(m[j].LastActivity())
}

func (m RoomList) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m RoomList) Sort() {
	sort.Sort(m)
}

func (m RoomList) ByID(id string) Room {
	for _, r := range m {
		if r.ID() == id {
			return r
		}
	}

	return nil
}
