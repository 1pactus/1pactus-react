package model

import "github.com/1pactus/1pactus-react/proto/gen/go/api"

type GlobalState struct {
	TimeIndex         int64 `gorm:"primaryKey;uniqueIndex:idx_time_index;not null"`
	Stake             int64 `gorm:"not null"`
	Supply            int64 `gorm:"not null"`
	CirculatingSupply int64 `gorm:"not null"`
	Txs               int64 `gorm:"not null"`
	Blocks            int64 `gorm:"not null"`
	Fee               int64 `gorm:"not null"`

	ActiveValidator int64 `gorm:"not null"`
	ActiveAccount   int64 `gorm:"not null"`

	ActiveValidatorDict map[string]bool `gorm:"-:all"`
	ActiveAccountDict   map[string]bool `gorm:"-:all"`
}

func NewGlobalState() *GlobalState {
	return &GlobalState{
		ActiveValidatorDict: make(map[string]bool),
		ActiveAccountDict:   make(map[string]bool),
	}
}

func (g *GlobalState) Reset(timeIndex int64) {
	g.TimeIndex = timeIndex
	g.Blocks = 0
	g.Txs = 0
	g.Fee = 0
	clear(g.ActiveValidatorDict)
	clear(g.ActiveAccountDict)
}

func (g *GlobalState) CreateCommitCopied() *GlobalState {
	return &GlobalState{
		TimeIndex:         g.TimeIndex,
		Stake:             g.Stake,
		Supply:            g.Supply,
		CirculatingSupply: g.CirculatingSupply,
		Txs:               g.Txs,
		Blocks:            g.Blocks,
		Fee:               g.Fee,
		ActiveValidator:   int64(len(g.ActiveValidatorDict)),
		ActiveAccount:     int64(len(g.ActiveAccountDict)),
	}
}

func (g *GlobalState) ToProto() *api.NetworkStatusData {
	return &api.NetworkStatusData{
		TimeIndex:         uint32(g.TimeIndex),
		Stake:             g.Stake,
		Supply:            g.Supply,
		CirculatingSupply: g.CirculatingSupply,
		Txs:               g.Txs,
		Blocks:            g.Blocks,
		Fee:               g.Fee,
		ActiveValidator:   g.ActiveValidator,
		ActiveAccount:     g.ActiveAccount,
	}
}
