package main

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/mgutz/ansi"
)

type UI interface {
	Render(g *gocui.Gui, state *State) error
	Scroll(g *gocui.Gui, mult float64) error
	Input(g *gocui.Gui) (string, error)
}

type ui struct{}

func NewUI() UI {
	return &ui{}
}

type State struct {
	Messages   []Message
	Rooms      []Room
	ActiveRoom Room
}

func (*ui) Render(g *gocui.Gui, state *State) error {
	roomBarWidth := 30
	inputHeight := 2
	maxX, yMax := g.Size()
	if v, err := g.SetView("chat", roomBarWidth, 0, maxX, yMax); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		v.Wrap = true
		drawMessages(v, state.Messages)
	}

	if v, err := g.SetView("rooms", 0, 0, roomBarWidth, yMax); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		drawRooms(v, state.Rooms, state.ActiveRoom)
	}

	if v, err := g.SetView("input", roomBarWidth, yMax-inputHeight, maxX, yMax); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		v.Editable = true
		v.Editor = &Editor{}
		if _, err := g.SetCurrentView(v.Name()); err != nil {
			return err
		}
	}

	return nil
}

func drawRooms(v *gocui.View, rooms []Room, active Room) {
	v.Clear()
	for _, r := range rooms {
		title := r.Title()
		if r == active {
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

func (*ui) Scroll(g *gocui.Gui, mult float64) error {
	v, err := g.View("chat")
	if err != nil {
		return err
	}

	_, h := v.Size()
	dy := int(float64(h) * mult)
	x, y := v.Origin()
	yNew := y + dy
	yMax := strings.Count(v.ViewBuffer(), "\n") - h + 1
	if yNew < 0 {
		yNew = 0
	} else if yNew > yMax {
		yNew = yMax
	}
	return v.SetOrigin(x, yNew)
}

func (*ui) Input(g *gocui.Gui) (string, error) {
	v, err := g.View("input")
	if err != nil {
		return "", err
	}
	return v.Buffer(), nil
}
