package menus

import (
	"lyn2n/i18n"
	"lyn2n/lib"
	"lyn2n/status"
	lyTheme "lyn2n/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// makeEditMenuSubItem 设置菜单的子菜单
func makeEditMenuSubItem(a fyne.App, w fyne.Window) []*fyne.MenuItem {
	// 创建配置菜单
	createCfg := fyne.NewMenuItem(i18n.Lang().CreateCfg, createCfg)
	// 删除配置菜单
	delcfg := fyne.NewMenuItem(i18n.Lang().DelCfg, delCfg)

	// 主题切换菜单
	themeMenu := fyne.NewMenuItem(i18n.Lang().Theme, nil)
	themeMenu.ChildMenu = fyne.NewMenu("", makeThemeMenuSubItem(a, w)...)

	// 隐藏至后台的菜单按钮
	hideInTrayMenu := fyne.NewMenuItem(i18n.Lang().HideInTrayMenu, func() {
		w.Hide()
		status.WindowsHideStatus.Set(true)
	})

	return []*fyne.MenuItem{createCfg, delcfg, themeMenu, hideInTrayMenu}
}

func createCfg() {
	cfgs := status.SupernodeConfigs.Get()
	cfgs = append(cfgs, &lib.Config{})
	status.SupernodeConfigs.Set(cfgs)
}

func delCfg() {
	todel := 0
	cfgs := status.SupernodeConfigs.Get()
	for i, e := range cfgs {
		if e.Selected {
			todel = i
			break
		}
	}
	cfgs = append(cfgs[:todel], cfgs[todel+1:]...)
	status.SupernodeConfigs.Set(cfgs)
}

// makeThemeMenuSubItem 用于创建主题切换菜单
func makeThemeMenuSubItem(a fyne.App, _ fyne.Window) []*fyne.MenuItem {
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
