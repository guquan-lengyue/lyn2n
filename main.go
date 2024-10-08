package main

import (
	"io"
	"log"
	"lyn2n/event"
	"lyn2n/menus"
	"lyn2n/statics"
	"lyn2n/views"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var VERSION = "20241005"

func main() {
	// 创建或打开日志文件
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)
	// 创建一个 MultiWriter，将日志输出到文件和终端
	log.Println("---------------------------------------------------------------------------")
	log.Println("---------------------VERSION:", VERSION, "------------------------------------")
	log.Println("---------------------------------------------------------------------------")
	a := app.NewWithID("guquanlengyue.n2n")
	a.SetIcon(statics.Icon)
	w := a.NewWindow("冷月N2N")
	menus.MakeTray(a, w)

	w.SetMainMenu(menus.Make(a, w))
	w.SetMaster()

	w.SetContent(views.MakeContent(a, w))
	w.Resize(fyne.NewSize(520, 520))
	w.ShowAndRun()

	event.CloseMainWindowsEvent.Triger(nil)
}
