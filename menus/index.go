package menus

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"lyn2n/i18n"
)

func Make(a fyne.App, w fyne.Window) *fyne.MainMenu {
	editMenu := fyne.NewMenu(i18n.Lang().Edit)
	editMenu.Items = makeEditMenuSubItem(a, w)
	return fyne.NewMainMenu(editMenu)
}

func MakeTray(a fyne.App, w fyne.Window) {
	showWindowFlag := true
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem(i18n.Lang().HideWindow, nil)
		h.Icon = theme.HomeIcon()
		menu := fyne.NewMenu("", h)
		h.Action = func() {
			if showWindowFlag {
				h.Label = i18n.Lang().ShowWindow
				w.Hide()
			} else {
				w.Show()
				h.Label = i18n.Lang().HideWindow
			}
			showWindowFlag = !showWindowFlag
			menu.Refresh()
		}
		desk.SetSystemTrayMenu(menu)
	}
}
