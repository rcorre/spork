package main

import (
	"log"
	"os"

	"github.com/jroimartin/gocui"
	"github.com/romana/rlog"
)

func main() {
	conf, err := LoadConfig("spork.yaml")
	if err != nil {
		log.Panicln(err)
	}
	rlog.Debug("Config:", conf)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true

	token, ok := os.LookupEnv("SPARK_TOKEN")
	if !ok {
		log.Panicln("SPARK_TOKEN must be set")
	}

	spark := NewSpark(conf.SparkURL, token)
	core := NewCore(spark)
	if err := core.LoadRooms(); err != nil {
		log.Panicln(err)
	}

	ui := NewUI(conf)
	manager := NewManager(core, ui)

	g.SetManager(manager)

	if err := manager.BindKeys(g, conf.Keys); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

//func listen(s Spark, m Manager, conf *Config) {
//	if err := s.Events.Register(); err != nil {
//		panic(err)
//	}
//	defer func() {
//		if err := s.Events.UnRegister(); err != nil {
//			rlog.Errorf("Failed to unregister websocket: %v", err)
//		} else {
//			rlog.Info("Device unregistered")
//		}
//	}()
//	msgChan, errChan, err := s.Events.Listen()
//	if err != nil {
//		panic(err)
//	}
//
//	interrupt := make(chan os.Signal, 1)
//	signal.Notify(interrupt, os.Interrupt)
//
//	for {
//		select {
//		case msg := <-msgChan:
//			rlog.Infof("msg: %s", msg)
//		case err := <-errChan:
//			rlog.Errorf("err: %v", err)
//		case <-interrupt:
//			return
//		}
//	}
//}
