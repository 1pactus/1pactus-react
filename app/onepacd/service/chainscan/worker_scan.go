package chainscan

import (
	"fmt"
	"sync"
	"time"

	"github.com/1pactus/1pactus-react/app/onepacd/constants"
	"github.com/1pactus/1pactus-react/app/onepacd/service/chainextract/chainreader"
	"github.com/1pactus/1pactus-react/app/onepacd/store"
	db "github.com/1pactus/1pactus-react/app/onepacd/store"
	"github.com/1pactus/1pactus-react/app/onepacd/store/model"
	"github.com/1pactus/1pactus-react/log"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type workerScan struct {
	log         log.ILogger
	grpcServers []string
	reader      chainreader.BlockchainReader
}

func newScanWorker(log log.ILogger, reader chainreader.BlockchainReader) *workerScan {
	p := &workerScan{
		log:    log,
		reader: reader,
	}

	return p
}

func (p *workerScan) getDailyStartHeight(height uint32) uint32 {
	if height <= 0 {
		return 0
	}

	blocksPerDay := uint32(8640)
	startHeight := ((height-1)/blocksPerDay)*blocksPerDay + 1

	return startHeight
}

func (p *workerScan) GetTimeIndex(timestamp uint32) int64 {
	t := time.Unix(int64(timestamp), 0).UTC()

	year, month, day := t.Date()

	dayStart := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	return dayStart.Unix()
}

func (p *workerScan) startCommit(wg *sync.WaitGroup) (chan *db.PgDBCommit, chan error) {
	wg.Add(1)

	commitChan := make(chan *db.PgDBCommit, 64)
	errorChan := make(chan error, 1)
	startTime := time.Now()

	go func() {
		for {
			commit := <-commitChan

			if commit == nil {
				p.log.Infof("commitChan closed")
				wg.Done()
				return
			}

			if err := store.Postgres.Commit(commit); err != nil {
				p.log.Errorf("commit failed: %v", err)
				errorChan <- err
			}

			processDuration := time.Since(startTime)

			p.log.Infof("commit height=%d/%d (%.2f%%) timeIndex=%d time=%v processDuration=%v",
				commit.GetHeight(), commit.GetLastBlockHeight(), float64(commit.GetHeight())/float64(commit.GetLastBlockHeight())*100, commit.GetTimeIndex(), time.Unix(int64(commit.GetTimeIndex()), 0).UTC(),
				processDuration)
		}
	}()

	return commitChan, errorChan
}

func (p *workerScan) FetchBlockchain(dieChan <-chan struct{}) error {
	defer p.log.Infof("FetchBlockchain exited")

	var commitWg sync.WaitGroup

	commitChan, commitErrChain := p.startCommit(&commitWg)

	var height int64
	var lastBlockHeight int64
	var globalState *model.GlobalState

	if state, err := store.Postgres.GetTopGlobalState(); err != nil {
		return fmt.Errorf("getTopGlobalState failed: %v", err)
	} else {
		if state != nil {
			globalState = state
		} else {
			globalState = model.NewGlobalState()
		}
	}

	topBlockInfo, err := store.Postgres.GetTopBlock()

	if err != nil {
		return fmt.Errorf("getTopBlock failed: %v", err)
	}

	if topBlockInfo != nil {
		height = topBlockInfo.Height
		//beginHeight = int(height)
		p.log.Infof("top block height: %v", height)
	}

	blockchainInfo, err := p.reader.GetBlockchainInfo()

	if err != nil {
		return fmt.Errorf("getBlockchainInfo failed: %v", err)
	}

	if blockchainInfo.IsPruned {
		return fmt.Errorf("blockchain is pruned")
	}

	lastBlockHeight = int64(blockchainInfo.LastBlockHeight)

	var lastTimeIndex int64
	txMerger := model.NewTxMerger()

	IsInitial := false

	group, _ := p.reader.CreateGroup(height+1, "pg_gatherer")

	defer group.Close()

	for {
		select {
		case <-dieChan:
			p.log.Warn("Context cancelled, stopping FetchBlockchain")
			return fmt.Errorf("cancelled")
		case err = <-commitErrChain:
			p.log.Errorf("commit error: %v", err)
			return err
		case block, ok := <-group.Read():
			if !ok {
				p.log.Infof("top height reached: %v, commit cancelled", height)
				commitChan <- nil // close commitChan
				commitWg.Wait()
				return nil
			}
			height = int64(block.Height)

			if height >= lastBlockHeight {
				group.Close()

				return nil
			}

			timeIndex := p.GetTimeIndex(block.BlockTime)

			if !IsInitial {
				IsInitial = true
				lastTimeIndex = timeIndex
				globalState.Reset(timeIndex)
			}

			// change day
			if timeIndex != lastTimeIndex {
				commitCtx := store.NewPgDBCommitContext(height, lastBlockHeight, lastTimeIndex, txMerger, globalState.CreateCommitCopied())
				commitChan <- commitCtx
				txMerger = model.NewTxMerger()

				lastTimeIndex = timeIndex
				globalState.Reset(timeIndex)
			}

			globalState.Txs += int64(len(block.Txs))
			globalState.Blocks += 1

			globalState.ActiveValidatorDict[block.Header.ProposerAddress] = true

			for _, tx := range block.Txs {
				globalState.Fee += tx.Fee
				switch tx.PayloadType {
				case pactus.PayloadType_PAYLOAD_TYPE_UNSPECIFIED:
				case pactus.PayloadType_PAYLOAD_TYPE_TRANSFER:
					globalState.ActiveAccountDict[tx.GetTransfer().Sender] = true

					if constants.IsMainnetReserveAccount(tx.GetTransfer().Sender) {
						globalState.Supply += tx.GetTransfer().Amount
						globalState.CirculatingSupply += tx.GetTransfer().Amount
					}

					if constants.IsMainnetReserveAccount(tx.GetTransfer().Receiver) {
						globalState.Supply -= tx.GetTransfer().Amount
						globalState.CirculatingSupply -= tx.GetTransfer().Amount
					}

					if constants.IsMainnetTeamHotAccount(tx.GetTransfer().Sender) {
						globalState.Supply += tx.GetTransfer().Amount
						globalState.CirculatingSupply += tx.GetTransfer().Amount
					}

					if constants.IsMainnetTeamHotAccount(tx.GetTransfer().Receiver) {
						globalState.Supply -= tx.GetTransfer().Amount
						globalState.CirculatingSupply -= tx.GetTransfer().Amount
					}

					if tx.GetTransfer().Sender == constants.Treasury {
						txMerger.AddReward(timeIndex, tx.GetTransfer().Receiver, tx.GetTransfer().Amount, block.Header.ProposerAddress)
					} else {
						txMerger.AddTransfer(timeIndex, tx.GetTransfer().Sender, tx.GetTransfer().Receiver, tx.GetTransfer().Amount, tx.Fee)
					}
				case pactus.PayloadType_PAYLOAD_TYPE_BOND:
					txMerger.AddBond(timeIndex, tx.GetBond().Sender, tx.GetBond().Receiver, tx.GetBond().Stake, tx.Fee)
					globalState.Stake += tx.GetBond().Stake
					globalState.CirculatingSupply -= tx.GetBond().Stake
				case pactus.PayloadType_PAYLOAD_TYPE_SORTITION:
				case pactus.PayloadType_PAYLOAD_TYPE_UNBOND:
					txMerger.AddUnbond(timeIndex, tx.GetUnbond().Validator, height, tx.GetId(), int64(block.BlockTime))
				case pactus.PayloadType_PAYLOAD_TYPE_WITHDRAW:
					txMerger.AddWithdraw(timeIndex, tx.GetWithdraw().ValidatorAddress, tx.GetWithdraw().AccountAddress, tx.GetWithdraw().Amount, tx.Fee)
					globalState.Stake -= tx.GetWithdraw().Amount
					globalState.CirculatingSupply += tx.GetWithdraw().Amount
				case pactus.PayloadType_PAYLOAD_TYPE_BATCH_TRANSFER:
					bt := tx.GetBatchTransfer()

					globalState.ActiveAccountDict[bt.Sender] = true

					for _, recipient := range bt.Recipients {
						if constants.IsMainnetReserveAccount(bt.Sender) {
							globalState.Supply += recipient.Amount
							globalState.CirculatingSupply += recipient.Amount
						}

						if constants.IsMainnetReserveAccount(recipient.Receiver) {
							globalState.Supply -= recipient.Amount
							globalState.CirculatingSupply -= recipient.Amount
						}

						if constants.IsMainnetTeamHotAccount(bt.Sender) {
							globalState.Supply += recipient.Amount
							globalState.CirculatingSupply += recipient.Amount
						}

						if constants.IsMainnetTeamHotAccount(recipient.Receiver) {
							globalState.Supply -= recipient.Amount
							globalState.CirculatingSupply -= recipient.Amount
						}

						if bt.Sender == constants.Treasury {
							txMerger.AddReward(timeIndex, recipient.Receiver, recipient.Amount, block.Header.ProposerAddress)
						} else {
							txMerger.AddTransfer(timeIndex, bt.Sender, recipient.Receiver, recipient.Amount, tx.Fee)
						}
					}
				}
			}
		}
	}
}
