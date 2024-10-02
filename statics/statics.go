package statics

import (
	_ "embed"
	"fyne.io/fyne/v2"
)

//go:embed icon.ico
var iconResource []byte
var Icon = fyne.NewStaticResource("icon", iconResource)
