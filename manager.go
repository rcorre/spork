package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/romana/rlog"
)

// Manager is the gocui.Manager for spork
type Manager interface {
	gocui.Manager

	BindKeys(g *gocui.Gui, keys map[string]string) error
	Handle(g *gocui.Gui, ev *Event) error

	// bindable commands
	NextRoom(g *gocui.Gui, _ *gocui.View) error
	PrevRoom(g *gocui.Gui, _ *gocui.View) error
	PageDown(g *gocui.Gui, _ *gocui.View) error
	PageUp(g *gocui.Gui, _ *gocui.View) error
	HalfPageDown(g *gocui.Gui, _ *gocui.View) error
	HalfPageUp(g *gocui.Gui, _ *gocui.View) error
	Send(g *gocui.Gui, _ *gocui.View) error
	Quit(g *gocui.Gui, _ *gocui.View) error
}

type manager struct {
	spark Spark
	view  UI
	core  *Core
	room  *Room
}

func NewManager(c *Core, v UI) Manager {
	return &manager{
		core: c,
		view: v,
		room: nil,
	}
}

func (m *manager) BindKeys(g *gocui.Gui, keys map[string]string) error {
	var keyMap = map[string]gocui.Key{
		"<c-a>":        gocui.KeyCtrlA,
		"<c-b>":        gocui.KeyCtrlB,
		"<c-c>":        gocui.KeyCtrlC,
		"<c-d>":        gocui.KeyCtrlD,
		"<c-e>":        gocui.KeyCtrlE,
		"<c-f>":        gocui.KeyCtrlF,
		"<c-g>":        gocui.KeyCtrlG,
		"<c-h>":        gocui.KeyCtrlH,
		"<c-i>":        gocui.KeyCtrlI,
		"<c-j>":        gocui.KeyCtrlJ,
		"<c-k>":        gocui.KeyCtrlK,
		"<c-l>":        gocui.KeyCtrlL,
		"<c-m>":        gocui.KeyCtrlM,
		"<c-n>":        gocui.KeyCtrlN,
		"<c-o>":        gocui.KeyCtrlO,
		"<c-p>":        gocui.KeyCtrlP,
		"<c-q>":        gocui.KeyCtrlQ,
		"<c-r>":        gocui.KeyCtrlR,
		"<c-s>":        gocui.KeyCtrlS,
		"<c-t>":        gocui.KeyCtrlT,
		"<c-u>":        gocui.KeyCtrlU,
		"<c-v>":        gocui.KeyCtrlV,
		"<c-w>":        gocui.KeyCtrlW,
		"<c-x>":        gocui.KeyCtrlX,
		"<c-y>":        gocui.KeyCtrlY,
		"<c-z>":        gocui.KeyCtrlZ,
		"<c-2>":        gocui.KeyCtrl2,
		"<c-3>":        gocui.KeyCtrl3,
		"<c-4>":        gocui.KeyCtrl4,
		"<c-5>":        gocui.KeyCtrl5,
		"<c-6>":        gocui.KeyCtrl6,
		"<c-7>":        gocui.KeyCtrl7,
		"<c-8>":        gocui.KeyCtrl8,
		"<c-~>":        gocui.KeyCtrlTilde,
		"<c-space>":    gocui.KeyCtrlSpace,
		"<backspace>":  gocui.KeyBackspace,
		"<tab>":        gocui.KeyTab,
		"<enter>":      gocui.KeyEnter,
		"<cr>":         gocui.KeyEnter,
		"<esc>":        gocui.KeyEsc,
		"<c-[>":        gocui.KeyCtrlLsqBracket,
		"<c-\\>":       gocui.KeyCtrlBackslash,
		"<c-]>":        gocui.KeyCtrlRsqBracket,
		"<c-/>":        gocui.KeyCtrlSlash,
		"<c-_>":        gocui.KeyCtrlUnderscore,
		"<space>":      gocui.KeySpace,
		"<backspace2>": gocui.KeyBackspace2,

		"<f1>":     gocui.KeyF1,
		"<f2>":     gocui.KeyF2,
		"<f3>":     gocui.KeyF3,
		"<f4>":     gocui.KeyF4,
		"<f5>":     gocui.KeyF5,
		"<f6>":     gocui.KeyF6,
		"<f7>":     gocui.KeyF7,
		"<f8>":     gocui.KeyF8,
		"<f9>":     gocui.KeyF9,
		"<f10>":    gocui.KeyF10,
		"<f11>":    gocui.KeyF11,
		"<f12>":    gocui.KeyF12,
		"<insert>": gocui.KeyInsert,
		"<delete>": gocui.KeyDelete,
		"<home>":   gocui.KeyHome,
		"<end>":    gocui.KeyEnd,
		"<pgup>":   gocui.KeyPgup,
		"<pgdn>":   gocui.KeyPgdn,
		"<up>":     gocui.KeyArrowUp,
		"<down>":   gocui.KeyArrowDown,
		"<left>":   gocui.KeyArrowLeft,
		"<right>":  gocui.KeyArrowRight,
	}

	cmdMap := map[string]func(*gocui.Gui, *gocui.View) error{
		"nextroom":     m.NextRoom,
		"prevroom":     m.PrevRoom,
		"pagedown":     m.PageDown,
		"pageup":       m.PageUp,
		"halfpagedown": m.HalfPageDown,
		"halfpageup":   m.HalfPageUp,
		"send":         m.Send,
		"quit":         m.Quit,
	}

	for keyName, cmdName := range keys {
		keyName := strings.ToLower(keyName)
		cmdName := strings.ToLower(cmdName)

		rlog.Debugf("Binding %q to %q", keyName, cmdName)

		key, ok := keyMap[keyName]
		if !ok {
			return fmt.Errorf("Unknown key %q. Should be one of %q.",
				cmdName, reflect.ValueOf(keyMap).MapKeys())
		}

		cmd, ok := cmdMap[cmdName]
		if !ok {
			return fmt.Errorf("Unknown command %q. Should be one of %q.",
				cmdName, reflect.ValueOf(cmdMap).MapKeys())
		}

		if err := g.SetKeybinding("", key, gocui.ModNone, cmd); err != nil {
			return err
		}
	}

	return nil
}

func (m *manager) Handle(g *gocui.Gui, ev *Event) error {
	rlog.Debugf("Processing spark event: %+v", ev)
	switch ev.Data.EventType {
	case "conversation.activity":
		activity := ev.Data.Activity
		fromID := activity.Actor.EntryUUID
		roomID := activity.Target.ID
		if activity.Verb == "acknowledge" {
			m.handleAcknowledge(g, roomID, fromID)
		} else if activity.Verb == "post" {
			msgID := activity.ID
			m.handleMessage(g, roomID, msgID)
		}
	case "status.start_typing":
		roomID := ev.Data.ConversationID
		fromID := ev.Data.Actor.EntryUUID
		m.handleStartTyping(g, roomID, fromID)
	case "status.stop_typing":
		roomID := ev.Data.ConversationID
		fromID := ev.Data.Actor.EntryUUID
		m.handleStopTyping(g, roomID, fromID)
	}
	return nil
}

func (m *manager) handleAcknowledge(g *gocui.Gui, roomID, personID string) error {
	return nil
}

func (m *manager) handleMessage(g *gocui.Gui, roomID, msgID string) error {
	room := m.core.Room(roomID)
	if room == nil {
		// TODO: try to load new room
		rlog.Warn("Could not find room %q for incoming message", roomID)
		return nil
	}

	return m.core.LoadMessages(roomID)
}

func (m *manager) handleStartTyping(g *gocui.Gui, roomID, personID string) error {
	return nil
}

func (m *manager) handleStopTyping(g *gocui.Gui, roomID, personID string) error {
	return nil
}

func (m *manager) updateRoom(g *gocui.Gui, r *Room) {
	g.Update(func(g *gocui.Gui) error {
		if err := m.core.LoadMessages(r.ID); err != nil {
			return err
		}

		if m.core.ActiveRoom == r {
			return m.view.Render(g, m.core)
		}
		return nil
	})
}

func (m *manager) NextRoom(g *gocui.Gui, _ *gocui.View) error {
	room := m.core.CycleRoom(1)
	m.updateRoom(g, room)
	return m.view.Render(g, m.core)
}

func (m *manager) PrevRoom(g *gocui.Gui, _ *gocui.View) error {
	room := m.core.CycleRoom(-1)
	m.updateRoom(g, room)
	return m.view.Render(g, m.core)
}

func (m *manager) PageUp(g *gocui.Gui, _ *gocui.View) error {
	return m.view.Scroll(g, -1)
}

func (m *manager) PageDown(g *gocui.Gui, _ *gocui.View) error {
	return m.view.Scroll(g, 1)
}

func (m *manager) HalfPageUp(g *gocui.Gui, _ *gocui.View) error {
	return m.view.Scroll(g, -1.0/2.0)
}

func (m *manager) HalfPageDown(g *gocui.Gui, _ *gocui.View) error {
	return m.view.Scroll(g, 1.0/2.0)
}

func (m *manager) Send(g *gocui.Gui, _ *gocui.View) error {
	text, err := m.view.Input(g)
	if err != nil {
		return err
	}

	id := m.core.ActiveRoom.ID
	if err := m.core.Send(text, id); err != nil {
		return err
	}

	m.updateRoom(g, m.core.ActiveRoom)
	return nil
}

func (m *manager) Layout(g *gocui.Gui) error {
	return m.view.Render(g, m.core)
}

func (m *manager) Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
