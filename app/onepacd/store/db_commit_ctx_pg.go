package store

import (
	"github.com/frimin/1pactus-react/app/onepacd/store/model"
)

type PgDBCommit struct {
	height          int64
	lastBlockHeight int64
	timeIndex       int64
	txMerger        *model.TxMerger
	globalState     *model.GlobalState
}

func NewPgDBCommitContext(height int64, lastBlockHeight int64, timeIndex int64, txMerger *model.TxMerger, globalState *model.GlobalState) *PgDBCommit {
	p := &PgDBCommit{
		height:          height,
		lastBlockHeight: lastBlockHeight,
		timeIndex:       timeIndex,
		txMerger:        txMerger,
		globalState:     globalState,
	}

	return p
}

func (c *PgDBCommit) GetTxMerger() *model.TxMerger {
	return c.txMerger
}

func (c *PgDBCommit) GetHeight() int64 {
	return c.height
}

func (c *PgDBCommit) GetLastBlockHeight() int64 {
	return c.lastBlockHeight
}

func (c *PgDBCommit) GetTimeIndex() int64 {
	return c.timeIndex
}

func (c *PgDBCommit) GetGlobalState() *model.GlobalState {
	return c.globalState
}
