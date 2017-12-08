package main

import "github.com/jroimartin/gocui"

type Editor struct{}

func (e *Editor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch key {
	case gocui.KeyEnter:
		return
	case gocui.KeyCtrlW:
		e.backwardsKillWord(v)
	default:
		gocui.DefaultEditor.Edit(v, key, ch, mod)
	}
}

func (e *Editor) backwardsKillWord(v *gocui.View) {
	x, _ := v.Cursor()
	b := v.Buffer()
	for x > 0 && b[x] != ' ' {
		v.EditDelete(true)
		x = x - 1
	}
}
