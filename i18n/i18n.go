package i18n

import (
	_ "embed"
	"encoding/json"
)

type Language struct {
	Edit           string `json:"edit"`
	Theme          string `json:"theme"`
	Dark           string `json:"dark"`
	Light          string `json:"light"`
	ShowWindow     string `json:"showWindow"`
	HideWindow     string `json:"hideWindow"`
	HideInTrayMenu string `json:"hideInTrayMenu"`
	Save           string `json:"save"`

	IpEntry        string `json:"ipEntry"`
	PortEntry      string `json:"portEntry"`
	RoomNameEntry  string `json:"roomName"`
	RoomKeyEntry   string `json:"roomKey"`
	EncryptedEntry string `json:"encryptedEntry"`
	StaticIpEntry  string `json:"StaticIpEntry"`
	ConnectText    string `json:"connectText"`
	DisconnectText string `json:"disconnectText"`

	ErrorInvalidIp        string `json:"errorInvalidIp"`
	ErrorInvalidPort      string `json:"errorInvalidPort"`
	ErrorRoomNameNotEmpty string `json:"errorRoomNameNotEmpty"`
}

//go:embed zh.json
var cnSource []byte

var cn *Language

func init() {
	cn = &Language{}
	err := json.Unmarshal(cnSource, cn)
	if err != nil {
		panic(err)
	}
}

func Lang() *Language {
	return cn
}
