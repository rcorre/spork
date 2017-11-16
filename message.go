package main

import (
	"log"
	"sort"
	"time"

	"github.com/rcorre/spork/spark"
)

type Message struct {
	Text   string
	Sender string
	Time   time.Time
}

// messageList implements sort.Interface to sort messages by time
type messageList []Message

func (m messageList) Len() int {
	return len(m)
}

func (m messageList) Less(i, j int) bool {
	return m[i].Time.Before(m[j].Time)
}

func (m messageList) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func LoadMessages(api *spark.Client, roomID string) ([]Message, error) {
	messages, err := api.Messages.List(roomID)
	if err != nil {
		return nil, err
	}

	ids := map[string]bool{}
	for _, msg := range messages {
		ids[msg.PersonID] = true
	}
	idList := []string{}
	for id, _ := range ids {
		idList = append(idList, id)
	}
	people, err := api.People.List(idList)
	if err != nil {
		return nil, err
	}

	out := make(messageList, len(messages))
	for i, msg := range messages {
		var name string
		for _, person := range people {
			if person.ID == msg.PersonID {
				name = person.DisplayName
				break
			}
		}

		time, err := time.Parse(time.RFC3339Nano, msg.Created)
		if err != nil {
			log.Printf("Failed to parse %s: %v", msg.Created, err)
			continue
		}
		out[i] = Message{
			Text:   msg.Text,
			Sender: name,
			Time:   time,
		}
	}

	sort.Sort(out)

	return out, nil
}
