package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"lyn2n/menus"
	"lyn2n/statics"
)

func main() {
	a := app.NewWithID("guquanlengyue.n2n")
	a.SetIcon(statics.Icon)
	w := a.NewWindow("冷月N2N")

	w.SetMainMenu(menus.Make(a, w))
	w.SetMaster()

	w.SetContent(widget.NewLabel("Hello World!"))
	w.Show()

	w.Resize(fyne.NewSize(520, 520))
	w.ShowAndRun()
}
