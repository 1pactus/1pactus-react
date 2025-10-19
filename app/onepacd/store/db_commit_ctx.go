package store

import (
	"github.com/1pactus/1pactus-react/app/onepacd/store/data"
)

type DBCommit struct {
	height          uint32
	lastBlockHeight uint32
	timeIndex       uint32
	txMerger        *TxMerger
	globalState     *data.GlobalStateData
}

func NewDBCommitContext(height uint32, lastBlockHeight uint32, timeIndex uint32, txMerger *TxMerger, globalState *data.GlobalStateData) *DBCommit {
	p := &DBCommit{
		height:          height,
		lastBlockHeight: lastBlockHeight,
		timeIndex:       timeIndex,
		txMerger:        txMerger,
		globalState:     globalState,
	}

	return p
}

func (c *DBCommit) GetTxMerger() *TxMerger {
	return c.txMerger
}

func (c *DBCommit) GetHeight() uint32 {
	return c.height
}

func (c *DBCommit) GetLastBlockHeight() uint32 {
	return c.lastBlockHeight
}

func (c *DBCommit) GetTimeIndex() uint32 {
	return c.timeIndex
}

func (c *DBCommit) GetGlobalState() *data.GlobalStateData {
	return c.globalState
}
