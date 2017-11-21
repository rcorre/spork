package main

import (
	"fmt"
	"text/tabwriter"

	"github.com/jroimartin/gocui"
)

type ChatView interface {
	Render(g *gocui.Gui, messages []Message) error
	Scroll(g *gocui.Gui, mult float64) error
}

type chatView struct{}

func NewChatView() ChatView {
	return &chatView{}
}

func (*chatView) Render(g *gocui.Gui, messages []Message) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("chat", 0, 0, maxX, maxY)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

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
