package main

import (
	"github.com/jroimartin/gocui"
	"github.com/rcorre/spork/spark"
)

const (
	CycleForward = 1
	CycleBackard = -1
)

type ChatController interface {
	NextRoom(g *gocui.Gui, _ *gocui.View) error
	PrevRoom(g *gocui.Gui, _ *gocui.View) error
	Layout(g *gocui.Gui) error
}

type chatController struct {
	spark   *spark.Client
	view    ChatView
	roomIdx int
	rooms   []spark.Room
}

func NewChatController(s *spark.Client, v ChatView) (ChatController, error) {
	rooms, err := s.Rooms.List()
	if err != nil {
		return nil, err
	}

	return &chatController{
		spark: s,
		view:  v,
		rooms: rooms,
	}, nil
}

func (c *chatController) NextRoom(g *gocui.Gui, _ *gocui.View) error {
	return c.cycleRoom(g, 1)
}

func (c *chatController) PrevRoom(g *gocui.Gui, _ *gocui.View) error {
	return c.cycleRoom(g, 1)
}

func (c *chatController) cycleRoom(g *gocui.Gui, direction int) error {
	c.roomIdx = (c.roomIdx + direction%len(c.rooms))
	room := c.rooms[c.roomIdx]
	messages, err := LoadMessages(c.spark, room.ID)
	if err != nil {
		return err
	}

	g.Update(func(g *gocui.Gui) error {
		return c.view.Render(g, messages)
	})
	return nil
}

func (c *chatController) Layout(g *gocui.Gui) error {
	messages, err := LoadMessages(c.spark, c.rooms[c.roomIdx].ID)
	if err != nil {
		return err
	}

	g.Update(func(g *gocui.Gui) error {
		return c.view.Render(g, messages)
	})
	return nil
}
