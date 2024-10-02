package views

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"log"
	"lyn2n/i18n"
	"lyn2n/lib"
	"net"
	"strconv"
)

func MakeContent(a fyne.App, w fyne.Window) fyne.CanvasObject {
	ipE := widget.NewEntry()
	ipE.Validator = func(s string) error {
		if net.ParseIP(s) == nil {
			return errors.New(i18n.Lang().ErrorInvalidIp)
		}
		return nil
	}

	portE := widget.NewEntry()
	portE.Validator = func(s string) error {
		errorInvalidPort := errors.New(i18n.Lang().ErrorInvalidPort)
		port, err := strconv.Atoi(s)
		if err != nil {
			return errorInvalidPort
		}
		if port < 0 || port > 65535 {
			return errorInvalidPort
		}
		return nil
	}

	roomNameE := widget.NewEntry()
	roomNameE.Validator = func(s string) error {
		if s == "" {
			return errors.New(i18n.Lang().ErrorRoomNameNotEmpty)
		}
		return nil
	}
	roomKeyE := widget.NewEntry()

	types := []string{
		"Twofish",
		"AES",
		"ChaCha20",
		"Speck-CTR",
	}
	encryptedE := widget.NewRadioGroup(types, func(s string) {
		log.Printf(s)
	})
	encryptedE.Horizontal = true
	encryptedE.Disable()

	items := []*widget.FormItem{
		{Text: i18n.Lang().IpEntry, Widget: ipE},
		{Text: i18n.Lang().PortEntry, Widget: portE},
		{Text: i18n.Lang().RoomNameEntry, Widget: roomNameE},
		{Text: i18n.Lang().RoomKeyEntry, Widget: roomKeyE},
		{Text: i18n.Lang().EncryptedEntry, Widget: encryptedE},
	}
	form := widget.NewForm(items...)
	form.SubmitText = i18n.Lang().ConnectText
	form.CancelText = i18n.Lang().ConnectText
	form.OnSubmit = func() {
		ip := ipE.Text
		port := portE.Text
		roomName := roomNameE.Text
		roomKey := roomKeyE.Text
		cmd := &lib.Command{
			Ip:       ip,
			Port:     port,
			RoomName: roomName,
			RoomKey:  roomKey,
			Encrypt:  encryptedE.Selected,
		}
		lib.Exec(cmd)
	}
	roomKeyE.OnChanged = func(s string) {
		if len(s) > 0 {
			encryptedE.Enable()
		} else {
			encryptedE.Disable()
		}
	}
	return form
}
