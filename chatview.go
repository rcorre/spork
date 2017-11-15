package main

import (
	"fmt"
	"text/tabwriter"

	"github.com/jroimartin/gocui"
)

type ChatView interface {
	Render(messages []Message) error
}

type chatView struct {
	chatPane *gocui.View
}

func NewChatView(g *gocui.Gui) (ChatView, error) {
	maxX, maxY := g.Size()
	chatPane, err := g.SetView("chat", 0, 0, maxX, maxY)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, err
	}

	return &chatView{
		chatPane: chatPane,
	}, nil
}

func (v *chatView) Render(messages []Message) error {
	w := new(tabwriter.Writer)
	w.Init(v.chatPane, 8, 8, 1, ' ', 0)
	for _, m := range messages {
		fmt.Fprintf(w, "%s\t| %s\n", m.Sender, m.Text)
	}
	return w.Flush()
}
