package views

import (
	"encoding/json"
	"log"
	"lyn2n/event"
	"lyn2n/lib"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var cfgs []*lib.Config

func loadCfgList(cfgs *[]*lib.Config) {
	file, err := os.OpenFile("cache.json", os.O_RDONLY, 0644)
	if err != nil {
		log.Println("Error opening cache.json", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(cfgs)
	if err != nil {
		log.Println("Error loading cache.json", err)
	}
}

func saveCfg() {
	file, err := os.OpenFile("cache.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("Error opening cache.json", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfgs)
	if err != nil {
		log.Println("Error loading cache.json", err)
	}
}

func MakeContent(a fyne.App, w fyne.Window) fyne.CanvasObject {
	loadCfgList(&cfgs)

	spnForm := SupernodeForm{}

	form := spnForm.MakeFormWidget(a, w)

	icon := widget.NewIcon(nil)
	hbox := container.NewHBox(icon, form)
	list := makeListWidget(&spnForm, icon)
	l := container.NewHSplit(list, container.NewCenter(hbox))
	spnForm.OnConnectFunc = func(c lib.Config) {
		saveCfg()
		l.Refresh()
	}
	return l
}

func makeListWidget(spnForm *SupernodeForm, icon *widget.Icon) fyne.CanvasObject {
	event.CloseMainWindowsEvent.Listen("OnCloseMinWindowsEventSaveConfigs", func(any) {
		saveCfg()
	})
	list := widget.NewList(
		func() int {
			return len(cfgs)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(nil), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(cfgs[id].ConfigName)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		cfg := cfgs[id]
		cfg.Selected = true
		spnForm.LoadCfg(cfg)
		icon.SetResource(nil)
	}

	list.OnUnselected = func(id widget.ListItemID) {
		cfgs[id].Selected = false
	}
	selectId := 0
	for i, cfg := range cfgs {
		if cfg.Selected {
			selectId = i
		}
	}
	list.Select(selectId)
	list.SetItemHeight(5, 50)
	list.SetItemHeight(6, 50)
	return list
}
