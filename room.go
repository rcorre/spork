package main

import (
	"github.com/rcorre/spork/spark"
)

// Room represents a spark room
type Room interface {
	Title() string
	Load() error
	Messages() []Message
}

type room struct {
	id       string
	title    string
	svc      spark.MessageService
	people   PersonCache
	messages MessageList
}

// NewRoom creates a Room wrapping a spark.Room
func NewRoom(src *spark.Room, svc spark.MessageService, people PersonCache) Room {
	return &room{
		id:     src.ID,
		title:  src.Title,
		svc:    svc,
		people: people,
	}
}

func (r *room) Title() string {
	return r.title
}

func (r *room) Messages() []Message {
	return r.messages
}

func (r *room) Load() error {
	messages, err := r.svc.List(r.id)
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
			Sender: sender,
			Time:   msg.Created,
		}
	}

	r.messages.Sort()
	return nil
}
