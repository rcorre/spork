package main

import "sort"

// Core contains the data that spork persists
type Core struct {
	Messages   map[string]MessageList
	People     map[string]*Person
	Rooms      RoomList
	Me         Person
	ActiveRoom *Room

	spark Spark
}

// NewCore returns a new Core
func NewCore(spark Spark) *Core {
	return &Core{
		spark: spark,
	}
}

// Room returns the room with the given id
func (c *Core) Room(id string) *Room {
	for _, r := range c.Rooms {
		if r.ID == id {
			return r
		}
	}

	return nil
}

func (c *Core) LoadRooms() error {
	rooms, err := c.spark.Rooms()
	if err != nil {
		return err
	}

	c.Rooms = RoomList(rooms)
	sort.Sort(c.Rooms)
	return nil
}

func (c *Core) LoadMessages(roomID string) error {
	messages, err := c.spark.Messages(roomID)
	if err != nil {
		return err
	}

	// load any unknown senders
	visited := map[string]bool{}
	ids := []string{}
	for _, msg := range messages {
		id := msg.PersonID
		if _, ok := c.People[id]; !ok {
			if _, ok := visited[id]; !ok {
				visited[msg.PersonID] = true
				ids = append(ids, id)
			}
		}
	}

	people, err := c.spark.People(ids)
	if err != nil {
		return err
	}

	for _, p := range people {
		c.People[p.ID] = p
	}

	msgList := MessageList(messages)
	sort.Sort(msgList)
	c.Messages[roomID] = append(c.Messages[roomID], msgList...)
	return nil
}

func (c *Core) Send(text, roomID string) error {
	msg, err := c.spark.Send(&Message{
		RoomID: roomID,
		Text:   text,
	})
	if err != nil {
		return err
	}

	c.Messages[roomID] = append(c.Messages[roomID], msg)
	return nil
}

// CycleRoom selects the next room in the given direction
func (c *Core) CycleRoom(direction int) *Room {
	curIdx := 0
	for i, r := range c.Rooms {
		if r == c.ActiveRoom {
			curIdx = i
			break
		}
	}
	idx := (curIdx + direction) % len(c.Rooms)
	c.ActiveRoom = c.Rooms[idx]
	return c.ActiveRoom
}

// RoomList implements sort.Interface to sort rooms by last activity
type RoomList []*Room

func (l RoomList) Len() int {
	return len(l)
}

func (l RoomList) Less(i, j int) bool {
	return l[i].LastActivity.After(*l[j].LastActivity)
}

func (l RoomList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l RoomList) ByID(id string) *Room {
	for _, r := range l {
		if r.ID == id {
			return r
		}
	}

	return nil
}

// MessageList implements sort.Interface to sort messages by time
type MessageList []*Message

func (l MessageList) Len() int {
	return len(l)
}

func (l MessageList) Less(i, j int) bool {
	return l[i].Created.Before(*l[j].Created)
}

func (l MessageList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
