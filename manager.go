package main

import (
	"github.com/jroimartin/gocui"
	"github.com/rcorre/spork/spark"
)

const (
	CycleForward = 1
	CycleBackard = -1
)

// Manager is the gocui.Manager for spork
type Manager interface {
	gocui.Manager

	NextRoom(g *gocui.Gui, _ *gocui.View) error
	PrevRoom(g *gocui.Gui, _ *gocui.View) error
	PageDown(g *gocui.Gui, _ *gocui.View) error
	PageUp(g *gocui.Gui, _ *gocui.View) error
	HalfPageDown(g *gocui.Gui, _ *gocui.View) error
	HalfPageUp(g *gocui.Gui, _ *gocui.View) error
}

type manager struct {
	spark   *spark.Client
	view    ChatView
	roomIdx int
	rooms   []Room
	people  PersonCache
}

func NewManager(s *spark.Client, v ChatView) (Manager, error) {
	roomList, err := s.Rooms.List()
	if err != nil {
		return nil, err
	}

	people := NewPersonCache(s.People)

	rooms := make([]Room, len(roomList))
	for i, r := range roomList {
		rooms[i] = NewRoom(&r, s.Messages, people)
		if err := rooms[i].Load(); err != nil {
			return nil, err
		}
	}

	return &manager{
		spark:  s,
		view:   v,
		rooms:  rooms,
		people: people,
	}, nil
}

func (m *manager) NextRoom(g *gocui.Gui, _ *gocui.View) error {
	return m.cycleRoom(g, 1)
}

func (m *manager) PrevRoom(g *gocui.Gui, _ *gocui.View) error {
	return m.cycleRoom(g, -1)
}

func (m *manager) cycleRoom(g *gocui.Gui, direction int) error {
	m.roomIdx = (m.roomIdx + direction%len(m.rooms))
	room := m.rooms[m.roomIdx]
	return m.view.Render(g, room.Messages())
}

func (m *manager) PageUp(g *gocui.Gui, _ *gocui.View) error {

	return m.view.Scroll(g, -1)
}

func (m *manager) PageDown(g *gocui.Gui, _ *gocui.View) error {
	return m.view.Scroll(g, 1)
}

func (m *manager) HalfPageUp(g *gocui.Gui, _ *gocui.View) error {

	return m.view.Scroll(g, -1.0/2.0)
}

func (m *manager) HalfPageDown(g *gocui.Gui, _ *gocui.View) error {
	return m.view.Scroll(g, 1.0/2.0)
}

func (m *manager) Layout(g *gocui.Gui) error {
	room := m.rooms[m.roomIdx]
	return m.view.Render(g, room.Messages())
}
