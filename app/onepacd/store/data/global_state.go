package data

import (
	"github.com/1pactus/1pactus-react/proto/gen/go/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GlobalState struct {
	TimeIndex         uint32 `bson:"time_index"`
	Stake             int64  `bson:"total_stake"`
	Supply            int64  `bson:"total_supply"`
	CirculatingSupply int64  `bson:"circulating_supply"`
	Txs               int64  `bson:"txs"`
	Blocks            int64  `bson:"blocks"`
	Fee               int64  `bson:"fee"`

	ActiveValidator map[string]bool
	ActiveAccount   map[string]bool
}

type GlobalStateData struct {
	ID                *primitive.ObjectID `bson:"_id,omitempty"`
	TimeIndex         uint32              `bson:"time_index"`
	Stake             int64               `bson:"total_stake"`
	Supply            int64               `bson:"total_supply"`
	CirculatingSupply int64               `bson:"circulating_supply"`
	Txs               int64               `bson:"txs"`
	Blocks            int64               `bson:"blocks"`
	Fee               int64               `bson:"fee"`

	ActiveValidatorCount int64 `bson:"active_validator_count"`
	ActiveAccountCount   int64 `bson:"active_account_count"`
}

func (g *GlobalStateData) ToProto() *api.NetworkStatusData {
	return &api.NetworkStatusData{
		TimeIndex:         g.TimeIndex,
		Stake:             g.Stake,
		Supply:            g.Supply,
		CirculatingSupply: g.CirculatingSupply,
		Txs:               g.Txs,
		Blocks:            g.Blocks,
		Fee:               g.Fee,
		ActiveValidator:   g.ActiveValidatorCount,
		ActiveAccount:     g.ActiveAccountCount,
	}
}

func NewGlobalStateData() *GlobalState {
	return &GlobalState{
		ActiveValidator: make(map[string]bool),
		ActiveAccount:   make(map[string]bool),
	}
}

func (g *GlobalState) Reset(timeIndex uint32) {
	g.TimeIndex = timeIndex
	g.Blocks = 0
	g.Txs = 0
	g.Fee = 0
	clear(g.ActiveValidator)
	clear(g.ActiveAccount)
}

func (g *GlobalState) CreateDBData() *GlobalStateData {
	if g == nil {
		return nil
	}

	copy := &GlobalStateData{
		TimeIndex:            g.TimeIndex,
		Stake:                g.Stake,
		Supply:               g.Supply,
		CirculatingSupply:    g.CirculatingSupply,
		Txs:                  g.Txs,
		Blocks:               g.Blocks,
		Fee:                  g.Fee,
		ActiveValidatorCount: int64(len(g.ActiveValidator)),
		ActiveAccountCount:   int64(len(g.ActiveAccount)),
		/*StakeChange:     g.StakeChange,
		SupplyChange:    g.SupplyChange,
		ValidatorsCount: g.ValidatorsCount,
		AccountsCount:   g.AccountsCount,*/
	}

	return copy
}
