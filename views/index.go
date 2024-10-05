package views

import (
	"encoding/json"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"log"
	"lyn2n/event"
	"lyn2n/i18n"
	"lyn2n/lib"
	"net"
	"os"
	"strconv"
)

var cmd *lib.Command

func MakeContent(a fyne.App, w fyne.Window) fyne.CanvasObject {
	cmd = &lib.Command{}
	load(cmd)
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
	staticIp := widget.NewEntry()
	go func() {
		for ip := range event.IpChange {
			staticIp.SetText(ip)
		}
	}()

	types := []string{
		"Twofish",
		"AES",
		"ChaCha20",
		"Speck-CTR",
	}
	encryptedE := widget.NewRadioGroup(types, nil)
	encryptedE.Horizontal = true
	encryptedE.Disable()

	items := []*widget.FormItem{
		{Text: i18n.Lang().IpEntry, Widget: ipE},
		{Text: i18n.Lang().PortEntry, Widget: portE},
		{Text: i18n.Lang().RoomNameEntry, Widget: roomNameE},
		{Text: i18n.Lang().RoomKeyEntry, Widget: roomKeyE},
		{Text: i18n.Lang().EncryptedEntry, Widget: encryptedE},
		{Text: i18n.Lang().StaticIpEntry, Widget: staticIp, HintText: i18n.Lang().FormHintMyIp},
	}
	form := widget.NewForm(items...)
	form.SubmitText = i18n.Lang().ConnectText
	form.CancelText = i18n.Lang().DisconnectText
	form.OnSubmit = onSubmit(ipE, portE, roomNameE, roomKeyE, encryptedE, staticIp)
	roomKeyE.OnChanged = func(s string) {
		if len(s) > 0 {
			encryptedE.SetSelected("AES")
			encryptedE.Enable()
		} else {
			encryptedE.SetSelected("")
			encryptedE.Disable()
		}
	}
	ipE.SetText(cmd.Ip)
	portE.SetText(cmd.Port)
	roomNameE.SetText(cmd.RoomName)
	roomKeyE.SetText(cmd.RoomKey)
	encryptedE.SetSelected(cmd.Encrypt)
	go func() {
		for {
			select {
			case <-event.N2NConnectedEvent:
				form.OnCancel = cmd.Kill
				form.OnSubmit = nil
			case <-event.N2NDisConnectedEvent:
				form.OnCancel = nil
				form.OnSubmit = onSubmit(ipE, portE, roomNameE, roomKeyE, encryptedE, staticIp)
			case <-event.N2NConnectedErr:
				form.OnCancel = cmd.Kill
				form.OnSubmit = nil
			}
			form.Refresh()
		}
	}()
	return form
}

func onSubmit(ipE *widget.Entry, portE *widget.Entry, roomNameE *widget.Entry, roomKeyE *widget.Entry, encryptedE *widget.RadioGroup, staticIp *widget.Entry) func() {
	return func() {
		cmd.Ip = ipE.Text
		cmd.Port = portE.Text
		cmd.RoomName = roomNameE.Text
		cmd.RoomKey = roomKeyE.Text
		cmd.Encrypt = encryptedE.Selected
		cmd.StaticIp = staticIp.Text
		go save(cmd)
		go cmd.Exec()
	}
}
func load(cmd *lib.Command) {
	file, err := os.OpenFile("cache.json", os.O_RDONLY, 0644)
	if err != nil {
		log.Println("Error opening cache.json", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(cmd)
	if err != nil {
		log.Println("Error loading cache.json", err)
	}
}

func save(cmd *lib.Command) {
	file, err := os.OpenFile("cache.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()
	if err != nil {
		log.Println("Error while opening cache.json", err)
		return
	}
	marshal, err := json.Marshal(cmd)
	if err != nil {
		log.Println("Error while marshalling cmd", err)
	}
	_, err = file.Write(marshal)
	if err != nil {
		log.Println("Error while writing to cache.json", err)
	}
}
