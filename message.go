package main

import (
	"sort"
	"time"
)

type Message struct {
	Text     string
	Markdown bool
	Sender   string
	Time     time.Time
}

// MessageList implements sort.Interface to sort messages by time
type MessageList []Message

func (m MessageList) Len() int {
	return len(m)
}

func (m MessageList) Less(i, j int) bool {
	return m[i].Time.Before(m[j].Time)
}

func (m MessageList) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m MessageList) Sort() {
	sort.Sort(m)
}
