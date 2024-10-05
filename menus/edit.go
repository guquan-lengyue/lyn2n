package menus

import (
	"lyn2n/i18n"
	"lyn2n/status"
	lyTheme "lyn2n/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func makeEditMenuSubItem(a fyne.App, w fyne.Window) []*fyne.MenuItem {
	themeMenu := fyne.NewMenuItem(i18n.Lang().Theme, nil)
	themeMenu.ChildMenu = fyne.NewMenu("", makeThemeMenuSubItem(a, w)...)

	hideInTrayMenu := fyne.NewMenuItem(i18n.Lang().HideInTrayMenu, func() {
		w.Hide()
		status.WindowsHideStatus = true
	})

	return []*fyne.MenuItem{themeMenu, hideInTrayMenu}
}

func makeThemeMenuSubItem(a fyne.App, w fyne.Window) []*fyne.MenuItem {
	themeDark := fyne.NewMenuItem(i18n.Lang().Dark, nil)
	themeLight := fyne.NewMenuItem(i18n.Lang().Light, nil)
	themeLight.Checked = true
	themeDark.Action = func() {
		a.Settings().SetTheme(&lyTheme.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantDark})
		themeDark.Checked = true
		themeLight.Checked = false
	}
	themeLight.Action = func() {
		a.Settings().SetTheme(&lyTheme.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantLight})
		themeLight.Checked = true
		themeDark.Checked = false
	}
	return []*fyne.MenuItem{themeLight, themeDark}
}
