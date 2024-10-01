package menus

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"lyn2n/i18n"
	lyTheme "lyn2n/theme"
)

func makeEditMenuSubItem(a fyne.App, w fyne.Window) []*fyne.MenuItem {
	themeMenu := fyne.NewMenuItem(i18n.Lang().Theme, nil)
	themeMenu.ChildMenu = fyne.NewMenu("", makeThemeMenuSubItem(a, w)...)

	hideInTrayMenu := fyne.NewMenuItem(i18n.Lang().HideInTrayMenu, func() {
		w.Hide()
	})
	return []*fyne.MenuItem{themeMenu, hideInTrayMenu}
}

func makeThemeMenuSubItem(a fyne.App, w fyne.Window) []*fyne.MenuItem {
	themeDark := fyne.NewMenuItem(i18n.Lang().Dark, func() {
		a.Settings().SetTheme(&lyTheme.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantDark})
	})
	themeLight := fyne.NewMenuItem(i18n.Lang().Light, func() {
		a.Settings().SetTheme(&lyTheme.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantLight})
	})
	return []*fyne.MenuItem{themeLight, themeDark}
}
