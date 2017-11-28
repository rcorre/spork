package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jroimartin/gocui"
	"github.com/rcorre/spork/spark"
	"github.com/romana/rlog"
)

func main() {
	runUI()
	//listen()
}

func listen() {
	conf, err := LoadConfig("spork.yaml")
	if err != nil {
		panic(err)
	}

	s, err := getSpark(conf)
	if err != nil {
		panic(err)
	}
	e := s.Events
	if err := e.Register(); err != nil {
		panic(err)
	}
	defer func() {
		if err := e.UnRegister(); err != nil {
			rlog.Errorf("Failed to unregister websocket: %v", err)
		} else {
			rlog.Info("Device unregistered")
		}
	}()
	msgChan, errChan, err := e.Listen()
	if err != nil {
		panic(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case msg := <-msgChan:
			rlog.Infof("msg: %s", msg)
		case err := <-errChan:
			rlog.Errorf("err: %v", err)
		case <-interrupt:
			return
		}
	}
}

func runUI() {
	conf, err := LoadConfig("spork.yaml")
	if err != nil {
		panic(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true

	s, err := getSpark(conf)
	if err != nil {
		panic(err)
	}

	manager, err := NewManager(s, NewChatView())
	if err != nil {
		panic(err)
	}

	g.SetManager(manager)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlJ, gocui.ModNone, manager.NextRoom); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlK, gocui.ModNone, manager.PrevRoom); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, manager.PageDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, manager.PageUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, manager.HalfPageDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlU, gocui.ModNone, manager.HalfPageUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, manager.Send); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func getSpark(conf *Config) (*spark.Client, error) {
	token, ok := os.LookupEnv("SPARK_TOKEN")
	if !ok {
		return nil, fmt.Errorf("SPARK_TOKEN must be set")
	}
	return spark.New(conf.SparkURL, conf.SparkDeviceURL, token), nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
