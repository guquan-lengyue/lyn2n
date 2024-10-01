package i18n

import (
	_ "embed"
	"encoding/json"
)

type Language struct {
	Edit  string `json:"edit"`
	Theme string `json:"theme"`
	Dark  string `json:"dark"`
	Light string `json:"light"`
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
