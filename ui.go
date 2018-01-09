package main

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/mgutz/ansi"
	"github.com/romana/rlog"
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
	statusHeight := 2
	maxX, yMax := g.Size()

	chatY2 := yMax - inputHeight - statusHeight
	if v, err := g.SetView("chat", roomBarWidth, 0, maxX, chatY2); err != nil && err != gocui.ErrUnknownView {
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

	// status needs to be drawn after the initial render so it can detect how
	// full the chat view is
	g.Update(func(g *gocui.Gui) error {
		statusY1 := yMax - inputHeight - statusHeight
		statusY2 := statusY1 + statusHeight
		rlog.Info(statusY1, statusY2)
		if v, err := g.SetView("status", roomBarWidth, statusY1, maxX, statusY2); err != nil && err != gocui.ErrUnknownView {
			return err
		} else {
			drawStatus(g, v)
		}
		return nil
	})

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
		if m.HTML != "" {
			fmt.Fprintln(v, HTMLtoText(m.HTML))
		} else {
			fmt.Fprintln(v, m.Text)
		}
	}
}

func drawStatus(g *gocui.Gui, v *gocui.View) {
	v.Clear()

	chatView, err := g.View("chat")
	if err != nil {
		rlog.Errorf("Failed to draw status: %v")
		return
	}

	_, y := chatView.Origin()
	_, h := chatView.Size()
	y = y + h
	yMax := strings.Count(chatView.ViewBuffer(), "\n")
	if yMax > 0 {
		fmt.Fprintf(v, "%d/%d (%d%%)", y, yMax, (y*100)/yMax)
	} else {
		fmt.Fprint(v, "0/0 (0%)")
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
	yMax := strings.Count(v.ViewBuffer(), "\n") - h
	if yMax < 0 {
		yMax = 0
	}
	if yNew < 0 {
		yNew = 0
	} else if yNew > yMax {
		yNew = yMax
	}
	rlog.Debugf("Scrolling from %d to %d", y, yNew)
	return v.SetOrigin(x, yNew)
}

func (*ui) Input(g *gocui.Gui) (string, error) {
	v, err := g.View("input")
	if err != nil {
		return "", err
	}
	return v.Buffer(), nil
}
