package menus

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"lyn2n/i18n"
	lyTheme "lyn2n/theme"
)

func Make(a fyne.App, w fyne.Window) *fyne.MainMenu {
	editMenu := fyne.NewMenu(i18n.Lang().Edit)
	for _, menu := range makeEditMenuSubItem(a, w) {
		editMenu.Items = append(editMenu.Items, menu)
	}
	return fyne.NewMainMenu(editMenu)
}

func makeEditMenuSubItem(a fyne.App, w fyne.Window) []*fyne.MenuItem {
	themeMenu := fyne.NewMenuItem(i18n.Lang().Theme, nil)

	themeDark := fyne.NewMenuItem(i18n.Lang().Dark, func() {
		a.Settings().SetTheme(&lyTheme.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantDark})
	})
	themeLight := fyne.NewMenuItem(i18n.Lang().Light, func() {
		a.Settings().SetTheme(&lyTheme.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantLight})
	})
	themeMenu.ChildMenu = fyne.NewMenu("", themeLight, themeDark)
	return []*fyne.MenuItem{themeMenu}
}
