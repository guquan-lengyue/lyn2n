package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"io"
	"log"
	"lyn2n/event"
	"lyn2n/menus"
	"lyn2n/statics"
	"lyn2n/views"
	"os"
)

func main() {
	// 创建或打开日志文件
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	// 创建一个 MultiWriter，将日志输出到文件和终端
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)
	a := app.NewWithID("guquanlengyue.n2n")
	a.SetIcon(statics.Icon)
	w := a.NewWindow("冷月N2N")
	menus.MakeTray(a, w)

	w.SetMainMenu(menus.Make(a, w))
	w.SetMaster()

	w.SetContent(views.MakeContent(a, w))
	w.SetOnClosed(func() {
		event.CloseMainWindowsEvent <- struct{}{}
	})

	w.Resize(fyne.NewSize(520, 520))
	w.ShowAndRun()
}
