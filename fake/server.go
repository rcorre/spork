package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rcorre/spork/spark"
)

type Room struct {
	Room     spark.Room
	Messages []spark.Message
}

var (
	rooms  []Room
	people []spark.Person
)

const testUserID = "07e1cc47-debd-4297-878f-976dea1b8081"

func init() {
	people = []spark.Person{{
		ID:          "0b38b012-b895-4fb4-93ab-8a2e116c3735",
		DisplayName: "Jane Doe",
	}, {
		ID:          "4da1d672-a098-4667-b97e-7a22111bfb8f",
		DisplayName: "Jim Bob",
	}, {
		ID:          testUserID,
		DisplayName: "Test User",
	}}

	rooms = []Room{{
		Room: spark.Room{
			ID:           "317eb0ef-b912-4f1f-be8d-c2b9db024305",
			Title:        "Room One",
			Type:         "group",
			IsLocked:     false,
			TeamId:       "8ba5c515-6dab-4a66-ab72-b229bb0ddbfc",
			LastActivity: time.Now(),
			Created:      time.Now(),
		},
		Messages: []spark.Message{{
			Text:     "Blah blah blah",
			PersonID: "4da1d672-a098-4667-b97e-7a22111bfb8f",
		}, {
			Text:     "whatever :/",
			PersonID: "0b38b012-b895-4fb4-93ab-8a2e116c3735",
		}, {
			Text:     "don't care -_-",
			PersonID: "0b38b012-b895-4fb4-93ab-8a2e116c3735",
		}},
	}, {
		Room: spark.Room{
			ID:           "d309325c-01bf-421b-bff5-a8f759a5e6d5",
			Title:        "Room Two",
			Type:         "group",
			IsLocked:     false,
			TeamId:       "8ba5c515-6dab-4a66-ab72-b229bb0ddbfc",
			LastActivity: time.Now(),
			Created:      time.Now(),
		},
		Messages: []spark.Message{{
			Text:     "Blah blah blah",
			PersonID: "4da1d672-a098-4667-b97e-7a22111bfb8f",
		}, {
			Text:     "whatever :/",
			PersonID: "0b38b012-b895-4fb4-93ab-8a2e116c3735",
		}, {
			Text:     "don't care -_-",
			PersonID: "0b38b012-b895-4fb4-93ab-8a2e116c3735",
		}},
	}}
}

func person(id string) *spark.Person {
	for _, p := range people {
		if p.ID == id {
			return &p
		}
	}
	return nil
}

func room(id string) *Room {
	for _, r := range rooms {
		if r.Room.ID == id {
			return &r
		}
	}
	return nil
}

func main() {
	http.HandleFunc("/v1/messages", handleMessages)
	http.HandleFunc("/v1/rooms", handleRooms)
	http.HandleFunc("/v1/people", handlePeople)
	http.HandleFunc("/v1/people/me", handleMe)

	log.Fatalln(http.ListenAndServe("localhost:3000", nil))
}

func writeItems(resp http.ResponseWriter, items interface{}) {
	write(resp, &struct {
		Items interface{}
	}{
		Items: items,
	})
}

func write(resp http.ResponseWriter, data interface{}) {
	if bytes, err := json.Marshal(data); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	} else if _, err := resp.Write(bytes); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}
}

func handleMessages(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		if roomID := req.URL.Query().Get("roomId"); roomID == "" {
			http.Error(resp, "roomId required", http.StatusBadRequest)
		} else if r := room(roomID); r == nil {
			http.Error(resp, "no such room", http.StatusBadRequest)
		} else {
			writeItems(resp, r.Messages)
		}
	case "POST":
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}

		var msg spark.Message

		if err := json.Unmarshal(body, &msg); err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}

		if msg.RoomID == "" {
			http.Error(resp, "msg must have roomID", http.StatusBadRequest)
			return
		}

		msg.Created = time.Now()
		r := room(msg.RoomID)
		if r == nil {
			http.Error(resp, "no such room", http.StatusBadRequest)
			return
		}

		r.Messages = append(r.Messages, msg)

		write(resp, msg)
	default:
		http.Error(resp, req.Method, http.StatusMethodNotAllowed)
	}
}

func handleRooms(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		writeItems(resp, rooms)
	default:
		http.Error(resp, req.Method, http.StatusMethodNotAllowed)
	}
}

func handlePeople(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(resp, req.Method, http.StatusMethodNotAllowed)
		return
	}

	ids := strings.Split(req.URL.Query().Get("id"), ",")
	ret := make([]spark.Person, len(ids))
	for i, id := range ids {
		if p := person(id); p != nil {
			ret[i] = *p
		} else {
			http.Error(resp, fmt.Sprintf("No person %q", id), http.StatusBadRequest)
			return
		}
	}
	writeItems(resp, ret)
}

func handleMe(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		write(resp, people[len(people)-1])
	default:
		http.Error(resp, req.Method, http.StatusMethodNotAllowed)
	}
}
