package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
	"github.com/rcorre/spork/spark"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("chat", 0, 0, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		token, ok := os.LookupEnv("SPARK_TOKEN")
		if !ok {
			return fmt.Errorf("SPARK_TOKEN must be set")
		}
		s := spark.New("", token)
		rooms, err := s.Rooms.List()
		if err != nil {
			panic(err)
		}
		messages, err := s.Messages.List(rooms[0].ID)
		if err != nil {
			panic(err)
		}
		for _, msg := range messages {
			fmt.Fprintln(v, msg.Text)
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
