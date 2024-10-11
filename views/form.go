package views

import (
	"errors"
	"lyn2n/event"
	"lyn2n/i18n"
	"lyn2n/lib"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type SupernodeForm struct {
	cmd *lib.Command

	ipE        *widget.Entry
	cfgNameE   *widget.Entry
	portE      *widget.Entry
	roomNameE  *widget.Entry
	roomKeyE   *widget.Entry
	encryptedE *widget.RadioGroup
	staticIp   *widget.Entry

	OnConnectFunc func(lib.Config)
}

func (f *SupernodeForm) MakeFormWidget(a fyne.App, w fyne.Window) fyne.CanvasObject {
	f.ipE = widget.NewEntry()
	f.ipE.Validator = func(s string) error {
		if len(s) == 0 {
			return errors.New(i18n.Lang().ErrorInvalidAddr)
		}
		return nil
	}

	f.cfgNameE = widget.NewEntry()
	f.cfgNameE.Validator = func(s string) error {
		if len(s) == 0 {
			return errors.New(i18n.Lang().ErrorInvalidCfgName)
		}
		return nil
	}

	f.portE = widget.NewEntry()
	f.portE.Validator = func(s string) error {
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

	f.roomNameE = widget.NewEntry()
	f.roomNameE.Validator = func(s string) error {
		if s == "" {
			return errors.New(i18n.Lang().ErrorRoomNameNotEmpty)
		}
		return nil
	}
	f.roomKeyE = widget.NewEntry()
	f.staticIp = widget.NewEntry()
	event.IpChange.Listen("ViewIpChangeEventHandle", func(ip string) {
		f.staticIp.SetText(ip)
	})
	types := []string{
		"Twofish",
		"AES",
		"ChaCha20",
		"Speck-CTR",
	}
	f.encryptedE = widget.NewRadioGroup(types, nil)
	f.encryptedE.Horizontal = true
	f.encryptedE.Disable()

	items := []*widget.FormItem{
		{Text: i18n.Lang().CfgName, Widget: f.cfgNameE},
		{Text: i18n.Lang().IpEntry, Widget: f.ipE},
		{Text: i18n.Lang().PortEntry, Widget: f.portE},
		{Text: i18n.Lang().RoomNameEntry, Widget: f.roomNameE},
		{Text: i18n.Lang().RoomKeyEntry, Widget: f.roomKeyE},
		{Text: i18n.Lang().EncryptedEntry, Widget: f.encryptedE},
		{Text: i18n.Lang().StaticIpEntry, Widget: f.staticIp, HintText: i18n.Lang().FormHintMyIp},
	}
	form := widget.NewForm(items...)
	form.SubmitText = i18n.Lang().ConnectText
	form.CancelText = i18n.Lang().DisconnectText
	form.OnSubmit = f.onSubmit
	f.roomKeyE.OnChanged = func(s string) {
		if len(s) > 0 {
			f.encryptedE.SetSelected("AES")
			f.encryptedE.Enable()
		} else {
			f.encryptedE.SetSelected("")
			f.encryptedE.Disable()
		}
	}
	event.N2NConnectedEvent.Listen("ViewsN2NConnectedEventHandle", func(any) {
		form.OnCancel = f.cmd.Stop
		form.OnSubmit = nil
		form.Refresh()
	})
	event.N2NDisConnectedEvent.Listen("ViewsN2NDisConnectedEventHandle", func(any) {
		form.OnCancel = nil
		form.OnSubmit = f.onSubmit
		form.Refresh()
	})
	event.N2NConnectedErr.Listen("ViewsN2NConnectedErrHandle", func(a any) {
		form.OnCancel = f.cmd.Stop
		form.OnSubmit = nil
		f.cmd.Stop()
		f.cmd.Kill()
		form.Refresh()
	})
	return form
}

func (f *SupernodeForm) LoadCfg(cmd *lib.Config) {
	if f.cmd == nil {
		f.cmd = &lib.Command{}
	}
	f.cmd.Config = cmd

	f.ipE.SetText(cmd.Ip)
	f.cfgNameE.SetText(cmd.ConfigName)
	f.portE.SetText(cmd.Port)
	f.roomNameE.SetText(cmd.RoomName)
	f.roomKeyE.SetText(cmd.RoomKey)
	f.encryptedE.SetSelected(cmd.Encrypt)
}

func (f *SupernodeForm) onSubmit() {
	f.cmd.Ip = f.ipE.Text
	f.cmd.ConfigName = f.cfgNameE.Text
	f.cmd.Port = f.portE.Text
	f.cmd.RoomName = f.roomNameE.Text
	f.cmd.RoomKey = f.roomKeyE.Text
	f.cmd.Encrypt = f.encryptedE.Selected
	f.cmd.StaticIp = f.staticIp.Text

	go f.save(*f.cmd.Config)
	go f.cmd.Exec()
}

func (f *SupernodeForm) save(cmd lib.Config) {
	if f.OnConnectFunc != nil {
		f.OnConnectFunc(cmd)
	}
}
