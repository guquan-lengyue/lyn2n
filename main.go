package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"lyn2n/menus"
	"lyn2n/statics"
	"lyn2n/views"
)

func main() {
	a := app.NewWithID("guquanlengyue.n2n")
	a.SetIcon(statics.Icon)
	w := a.NewWindow("冷月N2N")
	menus.MakeTray(a, w)

	w.SetMainMenu(menus.Make(a, w))
	w.SetMaster()

	w.SetContent(views.MakeContent(a, w))

	w.Resize(fyne.NewSize(520, 520))
	w.ShowAndRun()
}
