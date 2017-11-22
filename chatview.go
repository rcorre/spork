package main

import (
	"fmt"
	"text/tabwriter"

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
	maxX, maxY := g.Size()
	chatView, err := g.SetView("chat", roomBarWidth, 0, maxX, maxY)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err := drawMessages(chatView, state.Messages); err != nil {
		return err
	}

	roomView, err := g.SetView("rooms", 0, 0, roomBarWidth, maxY)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	drawRooms(roomView, state.Rooms, state.RoomIdx)

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

func drawMessages(v *gocui.View, messages []Message) error {
	w := new(tabwriter.Writer)
	v.Clear()
	w.Init(v, 8, 8, 1, ' ', 0)
	for _, m := range messages {
		fmt.Fprintf(w, "%s\t| %s\t| %s\n", m.Time, m.Sender, m.Text)
	}
	return w.Flush()
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
