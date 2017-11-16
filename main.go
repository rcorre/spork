package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
	"github.com/rcorre/spork/spark"
)

func main() {
	s, err := getSpark()
	if err != nil {
		panic(err)
	}
	e := s.Events
	if err := e.Register(); err != nil {
		panic(err)
	}
	defer func() {
		if err := e.UnRegister(); err != nil {
			log.Printf("Failed to unregister websocket: %v", err)
		} else {
			log.Printf("Device unregistered")
		}
	}()
	if err := e.Listen(); err != nil {
		panic(err)
	}
	// runUI()
}

func runUI() {
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

func getSpark() (*spark.Client, error) {
	token, ok := os.LookupEnv("SPARK_TOKEN")
	if !ok {
		return nil, fmt.Errorf("SPARK_TOKEN must be set")
	}
	return spark.New("", token), nil
}

func layout(g *gocui.Gui) error {
	s, err := getSpark()
	if err != nil {
		return err
	}

	rooms, err := s.Rooms.List()
	if err != nil {
		panic(err)
	}
	messages, err := LoadMessages(s, rooms[0].ID)
	if err != nil {
		return err
	}

	chatView, err := NewChatView(g)
	if err != nil {
		return err
	}
	return chatView.Render(messages)
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
