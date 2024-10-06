package menus

import (
	"lyn2n/i18n"
	"lyn2n/status"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

func Make(a fyne.App, w fyne.Window) *fyne.MainMenu {
	editMenu := fyne.NewMenu(i18n.Lang().Edit)
	editMenu.Items = makeEditMenuSubItem(a, w)
	return fyne.NewMainMenu(editMenu)
}

func MakeTray(a fyne.App, w fyne.Window) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem(i18n.Lang().HideWindow, nil)
		h.Icon = theme.HomeIcon()
		menu := fyne.NewMenu("", h)
		h.Action = func() {
			status.WindowsHideStatus.Set(!status.WindowsHideStatus.Get())
		}
		status.WindowsHideStatus.Listen("TrayHandleWindowsHideStatus", func(b bool) {
			if b {
				h.Label = i18n.Lang().ShowWindow
				w.Hide()
			} else {
				w.Show()
				h.Label = i18n.Lang().HideWindow
			}
			menu.Refresh()
		})
		desk.SetSystemTrayMenu(menu)
	}
}
