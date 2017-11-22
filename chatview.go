package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/mgutz/ansi"
)

type ChatView interface {
	Render(g *gocui.Gui, state *State) error
	Scroll(g *gocui.Gui, mult float64) error
}

type chatView struct{}

func NewChatView() ChatView {
	return &chatView{}
}

type State struct {
	Messages []Message
	Rooms    []Room
	RoomIdx  int
}

func (*chatView) Render(g *gocui.Gui, state *State) error {
	roomBarWidth := 30
	inputHeight := 2
	maxX, maxY := g.Size()
	if v, err := g.SetView("chat", roomBarWidth, 0, maxX, maxY); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		drawMessages(v, state.Messages)
	}

	if v, err := g.SetView("rooms", 0, 0, roomBarWidth, maxY); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		drawRooms(v, state.Rooms, state.RoomIdx)
	}

	if v, err := g.SetView("input", roomBarWidth, maxY-inputHeight, maxX, maxY); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		v.Editable = true
		g.SetCurrentView(v.Name())
	}

	return nil
}

func drawRooms(v *gocui.View, rooms []Room, current int) {
	v.Clear()
	for i, r := range rooms {
		title := r.Title()
		if i == current {
			title = ansi.Color(r.Title(), "white+b")
		}
		fmt.Fprintf(v, "%s\n", title)
	}
}

func drawMessages(v *gocui.View, messages []Message) {
	v.Clear()
	var curSender string
	for _, m := range messages {
		if m.Sender != curSender {
			curSender = m.Sender
			sender := ansi.Color(m.Sender, "white+b")
			fmt.Fprintf(v, "\n--- %s (%s)  ---\n", sender, m.Time)
		}
		fmt.Fprintln(v, m.Text)
	}
}

func (*chatView) Scroll(g *gocui.Gui, mult float64) error {
	v, err := g.View("chat")
	if err != nil {
		return err
	}

	_, h := v.Size()
	dy := int(float64(h) * mult)
	x, y := v.Origin()
	if y+dy >= 0 {
		return v.SetOrigin(x, y+dy)
	}
	return nil

}
