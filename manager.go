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
	Send(g *gocui.Gui, _ *gocui.View) error
}

type manager struct {
	spark      *spark.Client
	view       UI
	activeRoom Room
	rooms      []Room
	people     PersonCache
}

func NewManager(s *spark.Client, v UI) (Manager, error) {
	roomList, err := s.Rooms.List()
	if err != nil {
		return nil, err
	}

	people, err := NewPersonCache(s.People)
	if err != nil {
		return nil, err
	}

	rooms := make([]Room, len(roomList))
	for i, r := range roomList {
		rooms[i] = NewRoom(&r, s.Messages, people)
	}
	return &manager{
		spark:      s,
		view:       v,
		rooms:      rooms,
		activeRoom: rooms[0],
		people:     people,
	}, nil
}

func (m *manager) updateRoom(g *gocui.Gui, r Room) {
	g.Update(func(g *gocui.Gui) error {
		if err := r.Load(); err != nil {
			return err
		}

		if m.activeRoom == r {
			return m.view.Render(g, m.state())
		}
		return nil
	})
}

func (m *manager) state() *State {
	return &State{
		Messages:   m.activeRoom.Messages(),
		Rooms:      m.rooms,
		ActiveRoom: m.activeRoom,
	}
}

func (m *manager) NextRoom(g *gocui.Gui, _ *gocui.View) error {
	room, err := m.cycleRoom(g, 1)
	if err != nil {
		return err
	}
	m.updateRoom(g, room)
	return nil
}

func (m *manager) PrevRoom(g *gocui.Gui, _ *gocui.View) error {
	room, err := m.cycleRoom(g, -1)
	if err != nil {
		return err
	}
	m.updateRoom(g, room)
	return nil
}

func (m *manager) cycleRoom(g *gocui.Gui, direction int) (Room, error) {
	curIdx := 0
	for i, r := range m.rooms {
		if r == m.activeRoom {
			curIdx = i
			break
		}
	}
	idx := (curIdx + direction) % len(m.rooms)
	m.activeRoom = m.rooms[idx]
	err := m.view.Render(g, m.state())
	return m.activeRoom, err
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

func (m *manager) Send(g *gocui.Gui, _ *gocui.View) error {
	text, err := m.view.Input(g)
	if err != nil {
		return err
	}

	if err := m.activeRoom.Send(text); err != nil {
		return err
	}

	m.updateRoom(g, m.activeRoom)
	return nil
}

func (m *manager) Layout(g *gocui.Gui) error {
	return m.view.Render(g, m.state())
}
