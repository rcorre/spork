package main

import (
	"github.com/rcorre/spork/spark"
)

type Message struct {
	Text   string
	Sender string
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

	out := make([]Message, len(messages))
	for i, msg := range messages {
		var name string
		for _, person := range people {
			if person.ID == msg.PersonID {
				name = person.DisplayName
				break
			}
		}
		out[i] = Message{
			Text:   msg.Text,
			Sender: name,
		}
	}

	return out, nil
}
