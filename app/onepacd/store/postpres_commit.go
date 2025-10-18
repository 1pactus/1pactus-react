package store

import (
	"fmt"
	"sync"

	"github.com/frimin/1pactus-react/app/onepacd/store/model"
)

type PgCommitContext interface {
	GetTxMerger() *model.TxMerger
	GetHeight() int64
	GetTimeIndex() int64
	GetGlobalState() *model.GlobalState
}

func (c *postgresStore) Commit(commitContext PgCommitContext) error {
	updateFuncs := []struct {
		name string
		fn   func(commitContext PgCommitContext) error
	}{
		{"insertBlockData", c.insertBlockData},
		{"InsertGlobalState", c.insertGlobalState},
		/* {"updateSenderTransfers", c.updateSenderTransfers},
		{"updateReceiverTransfers", c.updateReceiverTransfers},
		{"updateRewardTransfers", c.updateRewardTransfers},
		{"updateSenderBond", c.updateSenderBond},
		{"updateReceiverBond", c.updateReceiverBond},
		{"updateUnbondTransfers", c.updateUnbondTransfers},
		{"updateWithdrawSender", c.updateWithdrawSender},
		{"updateWithdrawReceiver", c.updateWithdrawReceiver},
		{"updateAccountBalance", c.updateAccountBalance},
		{"updateAccountBalanceIndex", c.updateAccountBalanceIndex},
		{"updateValidatorStake", c.updateValidatorStake},*/
	}

	errChan := make(chan error, len(updateFuncs))
	var wg sync.WaitGroup

	for _, uf := range updateFuncs {
		wg.Add(1)
		go func(name string, fn func(PgCommitContext) error) {
			defer wg.Done()
			if err := fn(commitContext); err != nil {
				errChan <- fmt.Errorf("%s error: %w", name, err)
			}
		}(uf.name, uf.fn)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *postgresStore) insertBlockData(commitContext PgCommitContext) error {
	blockData := &model.Block{Height: commitContext.GetHeight(), TimeIndex: commitContext.GetTimeIndex()}
	if err := c.db.GetDB().Model(&model.Block{}).Create(blockData).Error; err != nil {
		return err
	}
	return nil
}

func (c *postgresStore) insertGlobalState(commitContext PgCommitContext) error {
	return c.InsertGlobalState(commitContext.GetGlobalState())
}
