package main

import (
	"fmt"
	"text/tabwriter"

	"github.com/jroimartin/gocui"
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
	if v, err := g.SetView("chat", roomBarWidth, 0, maxX, maxY); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		w := new(tabwriter.Writer)
		v.Clear()
		w.Init(v, 8, 8, 1, ' ', 0)
		for _, m := range state.Messages {
			fmt.Fprintf(w, "%s\t| %s\t| %s\n", m.Time, m.Sender, m.Text)
		}
		if err := w.Flush(); err != nil {
			return err
		}
	}
	if v, err := g.SetView("rooms", 0, 0, roomBarWidth, maxY); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		v.Clear()
		for i, r := range state.Rooms {
			if i == state.RoomIdx {
				fmt.Fprintf(v, "- %s\n", r.Title())
			} else {
				fmt.Fprintf(v, "%s\n", r.Title())
			}
		}
	}

	return nil
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
