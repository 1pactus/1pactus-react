package constants

import (
	_ "embed"
	"encoding/json"

	"github.com/1pactus/1pactus-react/log"
)

//go:embed accounts_genesis.json
var accountsGenesisBytes []byte

//go:embed accounts_foundation_pip43.json
var accountsFoundationPip43Bytes []byte

var accountsGenesis []*AccountsGenesis
var accountsFoundationPip43 []string

type AccountsGenesis struct {
	Address string `json:"address"`
	Balance int64  `json:"balance"`
}

func init() {
	if err := json.Unmarshal(accountsGenesisBytes, &accountsGenesis); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(accountsFoundationPip43Bytes, &accountsFoundationPip43); err != nil {
		panic(err)
	}

	log.Infof("Loaded %d genesis accounts", len(accountsGenesis))
}
