package store

import (
	"context"
	"fmt"
	"sync"

	"github.com/frimin/1pactus-react/app/onepacd/store/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TreasuryAddress = "000000000000000000000000000000000000000000"
)

type DbClient struct {
	ctx context.Context
	//dbName     string
	//client     *mongo.Client
	collection *DbCollection
	database   *mongo.Database
}

type DbCollection struct {
	//addresses                       *mongo.Collection
	global_state_index              *mongo.Collection
	address_transfer_sender_index   *mongo.Collection
	address_transfer_receiver_index *mongo.Collection
	address_transfer_reward_index   *mongo.Collection
	address_bond_sender_index       *mongo.Collection
	address_bond_receiver_index     *mongo.Collection
	address_unbond_index            *mongo.Collection
	address_withdraw_sender_index   *mongo.Collection
	address_withdraw_receiver_index *mongo.Collection

	account_balance       *mongo.Collection
	account_balance_index *mongo.Collection
	validator_stake       *mongo.Collection

	block *mongo.Collection
	//tx                              *mongo.Collection
}

func NewDBClient() *DbClient {
	ctx := context.Background()

	db := &DbClient{
		ctx:        ctx,
		collection: &DbCollection{},
	}

	return db
}

func (c *DbClient) Connect(database *mongo.Database) error {
	c.database = database

	/*client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))

	if err != nil {
		return err
	}*/

	//c.collection.addresses = client.Database(dbName).Collection("addresses")

	c.collection.block = c.database.Collection("block")
	//c.collection.tx = client.Database(dbName).Collection("tx")

	/*if err := data.CreateAddressDataIndex(c.collection.addresses); err != nil {
		return err
	}*/

	if err := c.createAddressIndex(); err != nil {
		return err
	}

	if err := data.CreateBlockDataIndex(c.collection.block); err != nil {
		return err
	}

	if err := c.createAddressStateIndex(); err != nil {
		return err
	}

	if err := c.createGlobalStateIndex(); err != nil {
		return err
	}

	/*if err := data.CreateTxDataIndex(c.collection.tx); err != nil {
		return err
	}*/

	return nil
}

/*
func (c *dbClient) insertBlock(block *pactus.GetBlockResponse) error {
	// Start a session
	session, err := c.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(c.ctx)

	// Start a transaction
	_, err = session.WithTransaction(c.ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Insert block data
		blockData := &data.BlockData{
			Height: block.Height,
		}

		_, err := c.collection.block.InsertOne(sessCtx, blockData)
		if err != nil {
			return nil, err
		}

		// Prepare transaction data for bulk insert
		var txDocs []interface{}
		for _, tx := range block.Txs {
			newTxData := &data.TransactionData{
				Height: block.Height,
				TxId:   tx.Id,
				Value:  tx.Value,
				Fee:    tx.Fee,
				Memo:   tx.Memo,
			}

			switch tx.PayloadType {
			case pactus.PayloadType_UNKNOWN:
			case pactus.PayloadType_TRANSFER_PAYLOAD:
				newTxData.Transfer = &data.TransferPayloadData{
					Sender: tx.GetTransfer().Sender,
					Amount: tx.GetTransfer().Amount,
				}
				sender := tx.GetTransfer().Sender

				if sender != TreasuryAddress {
					r := tx.GetTransfer().Receiver
					log.Print(r)
				}
			case pactus.PayloadType_BOND_PAYLOAD:
			case pactus.PayloadType_SORTITION_PAYLOAD:
			case pactus.PayloadType_UNBOND_PAYLOAD:
			case pactus.PayloadType_WITHDRAW_PAYLOAD:
			}

			txDocs = append(txDocs, newTxData)
		}

		//pactus.AddressType.

		// Bulk insert all transaction data
		if len(txDocs) > 0 {
			_, err := c.collection.tx.InsertMany(sessCtx, txDocs)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return err
}*/

type CommitContext interface {
	GetTxMerger() *TxMerger
	GetHeight() uint32
	GetTimeIndex() uint32
	GetGlobalState() *data.GlobalStateData
}

func (c *DbClient) Commit(commitContext CommitContext) error {
	updateFuncs := []struct {
		name string
		fn   func(commitContext CommitContext) error
	}{
		{"insertBlockData", c.insertBlockData},
		{"InsertGlobalState", c.InsertGlobalState},
		{"updateSenderTransfers", c.updateSenderTransfers},
		{"updateReceiverTransfers", c.updateReceiverTransfers},
		{"updateRewardTransfers", c.updateRewardTransfers},
		{"updateSenderBond", c.updateSenderBond},
		{"updateReceiverBond", c.updateReceiverBond},
		{"updateUnbondTransfers", c.updateUnbondTransfers},
		{"updateWithdrawSender", c.updateWithdrawSender},
		{"updateWithdrawReceiver", c.updateWithdrawReceiver},
		{"updateAccountBalance", c.updateAccountBalance},
		{"updateAccountBalanceIndex", c.updateAccountBalanceIndex},
		{"updateValidatorStake", c.updateValidatorStake},
	}

	errChan := make(chan error, len(updateFuncs))
	var wg sync.WaitGroup

	for _, uf := range updateFuncs {
		wg.Add(1)
		go func(name string, fn func(CommitContext) error) {
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

func (c *DbClient) CleanupBlockCollection() error {
	opts := options.FindOne().SetSort(bson.D{{Key: "height", Value: -1}})
	var latestBlock data.BlockData
	err := c.collection.block.FindOne(context.Background(), bson.D{}, opts).Decode(&latestBlock)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return fmt.Errorf("get latest block failed: %v", err)
	}

	filter := bson.D{{Key: "height", Value: bson.D{{Key: "$lt", Value: latestBlock.Height}}}}
	_, err = c.collection.block.DeleteMany(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("delete block failed: %v", err)
	}

	return nil
}

func (c *DbClient) insertBlockData(commitContext CommitContext) error {
	blockData := &data.BlockData{Height: commitContext.GetHeight(), TimeIndex: commitContext.GetTimeIndex()}
	if _, err := c.collection.block.InsertOne(c.ctx, blockData); err != nil {
		return err
	}
	return nil
}

func (c *DbClient) GetTopBlock() (*data.BlockData, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "height", Value: -1}})

	var result data.BlockData
	err := c.collection.block.FindOne(c.ctx, bson.D{}, opts).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *DbClient) GetTopGlobalState() (*data.GlobalState, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "time_index", Value: -1}})

	result := data.NewGlobalStateData()
	err := c.collection.global_state_index.FindOne(c.ctx, bson.D{}, opts).Decode(result)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *DbClient) InsertGlobalState(commitContext CommitContext) error {
	_, err := c.collection.global_state_index.InsertOne(c.ctx, commitContext.GetGlobalState())
	if err != nil {
		return err
	}
	return nil
}
