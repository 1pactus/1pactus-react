package gather

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/frimin/1pactus-react/app/onepacd/store"
	db "github.com/frimin/1pactus-react/app/onepacd/store"
	"github.com/frimin/1pactus-react/app/onepacd/store/data"
	"github.com/frimin/1pactus-react/log"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

const (
	Treasury = "000000000000000000000000000000000000000000"
)

type ChainGather struct {
	db   *db.DbClient
	grpc *grpcClient
	log  log.ILogger
}

func NewChainGather(log log.ILogger, grpcServers []string) *ChainGather {
	p := &ChainGather{
		db:   db.NewDBClient(),
		grpc: newGrpcClient(time.Second*10, grpcServers),
		log:  log,
	}

	return p
}

func (p *ChainGather) Connect() error {
	p.db = store.Mongo.GetDBAdapter()

	if p.db == nil {
		return fmt.Errorf("db is not initialized")
	}

	err := p.grpc.connect()

	if err != nil {
		return err
	}

	return nil
}

func getDailyStartHeight(height uint32) uint32 {
	if height <= 0 {
		return 0
	}

	blocksPerDay := uint32(8640)
	startHeight := ((height-1)/blocksPerDay)*blocksPerDay + 1

	return startHeight
}

func getDayStartTimestamp(timestamp uint32) uint32 {
	// 将 uint32 转换为 int64，因为 time.Unix 需要 int64 类型
	t := time.Unix(int64(timestamp), 0).UTC()

	// 获取当天的年、月、日
	year, month, day := t.Date()

	// 创建当天 UTC+0 时间的起始时间
	dayStart := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	// 将时间转换回 Unix 时间戳，并转为 uint32
	return uint32(dayStart.Unix())
}

func (p *ChainGather) startCommit(wg *sync.WaitGroup) (chan *db.DBCommit, chan error) {
	wg.Add(1)

	commitChan := make(chan *db.DBCommit, 64)
	errorChan := make(chan error, 1)

	go func() {
		for {
			commit := <-commitChan

			if commit == nil {
				p.log.Infof("commitChan closed")
				wg.Done()
				return
			}

			if err := p.db.Commit(commit); err != nil {
				p.log.Errorf("commit failed: %v", err)
				errorChan <- err
			}

			p.log.Infof("commit height=%d/%d (%.2f%%) timeIndex=%d time=%v",
				commit.GetHeight(), commit.GetLastBlockHeight(), float64(commit.GetHeight())/float64(commit.GetLastBlockHeight())*100, commit.GetTimeIndex(), time.Unix(int64(commit.GetTimeIndex()), 0).UTC())
		}
	}()

	return commitChan, errorChan
}

func (p *ChainGather) FetchBlockchain(ctx context.Context) error {
	_, err := p.grpc.getBlockchainInfo()
	if err != nil {
		return err
	}

	var commitWg sync.WaitGroup

	commitChan, commitErrChain := p.startCommit(&commitWg)

	//beginHeight := 1
	var height uint32
	var lastBlockHeight uint32
	var globalState *data.GlobalState

	if state, err := p.db.GetTopGlobalState(); err != nil {
		return fmt.Errorf("getTopGlobalState failed: %v", err)
	} else {
		if state != nil {
			globalState = state
		} else {
			globalState = data.NewGlobalStateData()
		}
	}

	topBlockInfo, err := p.db.GetTopBlock()

	if err != nil {
		return fmt.Errorf("getTopBlock failed: %v", err)
	}
	if topBlockInfo != nil {
		height = topBlockInfo.Height
		//beginHeight = int(height)
		p.log.Infof("topBlockInfo.Height=%v", height)
	}

	blockchainInfo, err := p.grpc.getBlockchainInfo()

	if err != nil {
		return fmt.Errorf("getBlockchainInfo failed: %v", err)
	}

	if blockchainInfo.IsPruned {
		return fmt.Errorf("blockchain is pruned")
	}

	lastBlockHeight = blockchainInfo.LastBlockHeight

	var lastTimeIndex uint32
	txMerger := db.NewTxMerger()

	IsInitial := false

	for {
		select {
		case <-ctx.Done():
			p.log.Warn("Context cancelled, stopping FetchBlockchain")
			return ctx.Err()
		case err = <-commitErrChain:
			p.log.Errorf("commit error: %v", err)
			return err
		default:
			height++

			if height >= lastBlockHeight {
				p.log.Infof("top height reached: %v", height)

				commitChan <- nil // close commitChan
				commitWg.Wait()

				return nil
			}

			block, err := p.grpc.getBlock(height, pactus.BlockVerbosity_BLOCK_VERBOSITY_TRANSACTIONS)

			if err != nil {
				p.log.Errorf("getBlock failed: %v", err.Error())

				return err
			}

			timeIndex := getDayStartTimestamp(block.BlockTime)

			if !IsInitial {
				IsInitial = true
				lastTimeIndex = timeIndex
				globalState.Reset(timeIndex)
			}

			// change day
			if timeIndex != lastTimeIndex {
				commitCtx := db.NewDBCommitContext(height, lastBlockHeight, lastTimeIndex, txMerger, globalState.CreateDBData())
				commitChan <- commitCtx
				txMerger = db.NewTxMerger()

				lastTimeIndex = timeIndex
				globalState.Reset(timeIndex)
			}

			globalState.Txs += int64(len(block.Txs))
			globalState.Blocks += 1

			globalState.ActiveValidator[block.Header.ProposerAddress] = true

			for _, tx := range block.Txs {
				globalState.Fee += tx.Fee
				switch tx.PayloadType {
				case pactus.PayloadType_PAYLOAD_TYPE_UNSPECIFIED:
				case pactus.PayloadType_PAYLOAD_TYPE_TRANSFER:
					globalState.ActiveAccount[tx.GetTransfer().Sender] = true

					if p.db.IsMainnetReserveAccount(tx.GetTransfer().Sender) {
						globalState.Supply += tx.GetTransfer().Amount
						globalState.CirculatingSupply += tx.GetTransfer().Amount
					}

					if p.db.IsMainnetReserveAccount(tx.GetTransfer().Receiver) {
						globalState.Supply -= tx.GetTransfer().Amount
						globalState.CirculatingSupply -= tx.GetTransfer().Amount
					}

					if p.db.IsMainnetTeamHotAccount(tx.GetTransfer().Sender) {
						globalState.Supply += tx.GetTransfer().Amount
						globalState.CirculatingSupply += tx.GetTransfer().Amount
					}

					if p.db.IsMainnetTeamHotAccount(tx.GetTransfer().Receiver) {
						globalState.Supply -= tx.GetTransfer().Amount
						globalState.CirculatingSupply -= tx.GetTransfer().Amount
					}

					if tx.GetTransfer().Sender == Treasury {
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
					txMerger.AddUnbond(timeIndex, tx.GetUnbond().Validator, height, tx.GetId(), block.BlockTime)
				case pactus.PayloadType_PAYLOAD_TYPE_WITHDRAW:
					txMerger.AddWithdraw(timeIndex, tx.GetWithdraw().ValidatorAddress, tx.GetWithdraw().AccountAddress, tx.GetWithdraw().Amount, tx.Fee)
					globalState.Stake -= tx.GetWithdraw().Amount
					globalState.CirculatingSupply += tx.GetWithdraw().Amount
				case pactus.PayloadType_PAYLOAD_TYPE_BATCH_TRANSFER:
					bt := tx.GetBatchTransfer()

					globalState.ActiveAccount[bt.Sender] = true

					for _, recipient := range bt.Recipients {
						if p.db.IsMainnetReserveAccount(bt.Sender) {
							globalState.Supply += recipient.Amount
							globalState.CirculatingSupply += recipient.Amount
						}

						if p.db.IsMainnetReserveAccount(recipient.Receiver) {
							globalState.Supply -= recipient.Amount
							globalState.CirculatingSupply -= recipient.Amount
						}

						if p.db.IsMainnetTeamHotAccount(bt.Sender) {
							globalState.Supply += recipient.Amount
							globalState.CirculatingSupply += recipient.Amount
						}

						if p.db.IsMainnetTeamHotAccount(recipient.Receiver) {
							globalState.Supply -= recipient.Amount
							globalState.CirculatingSupply -= recipient.Amount
						}

						if bt.Sender == Treasury {
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
