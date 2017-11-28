package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/rcorre/spork/spark"
)

var (
	messages []spark.Message
	rooms    []spark.Room
	people   []spark.Person
)

func main() {
	http.HandleFunc("/v1/messages", handleMessages)
	http.HandleFunc("/v1/rooms", handleRooms)
	http.HandleFunc("/v1/people", handlePeople)

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
		writeItems(resp, messages)
	case "POST":
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
		}

		var msg spark.Message

		if err := json.Unmarshal(body, &msg); err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
		}

		msg.Created = time.Now()
		messages = append(messages, msg)

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
	switch req.Method {
	case "GET":
		writeItems(resp, people)
	default:
		http.Error(resp, req.Method, http.StatusMethodNotAllowed)
	}
}
