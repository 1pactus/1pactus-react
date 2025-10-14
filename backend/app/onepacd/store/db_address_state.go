package store

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createUniqueAddressIndex(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "address", Value: 1}},
		Options: options.Index().
			SetUnique(false),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err
	}

	return nil
}

func (c *DbClient) IsMainnetReserveAccount(account string) bool {
	switch account {
	case "000000000000000000000000000000000000000000",
		"pc1z2r0fmu8sg2ffa0tgrr08gnefcxl2kq7wvquf8z",
		"pc1zprhnvcsy3pthekdcu28cw8muw4f432hkwgfasv",
		"pc1znn2qxsugfrt7j4608zvtnxf8dnz8skrxguyf45",
		"pc1zs64vdggjcshumjwzaskhfn0j9gfpkvche3kxd3":
		return true
	default:
		return false
	}
}

func (c *DbClient) IsMainnetTeamHotAccount(account string) bool {
	switch account {
	// bootstarp reward account
	/*case "pc1zc7ndap6mx2znve365cknnmg20umtvxm50nmmlt",
	"pc1zp30eyll5vygs30x0j9mgpl7pj3mq9gakkuw87t",
	"pc1zvt3vhu9mhhq3lcuakz0gm00egz5fjf0zq4uzjd",
	"pc1zpjxwj4a5ssuh4vjgcfwwzd0z6zhlpj8ylnhdl8":
	return true*/
	case "pc1zuavu4sjcxcx9zsl8rlwwx0amnl94sp0el3u37g",
		"pc1zf0gyc4kxlfsvu64pheqzmk8r9eyzxqvxlk6s6t":
		return true
	default:
		return false
	}
}

func (c *DbClient) createAddressStateIndex() error {
	c.collection.account_balance = c.database.Collection("account_balance")
	c.collection.account_balance_index = c.database.Collection("account_balance_index")
	c.collection.validator_stake = c.database.Collection("validator_stake")

	if err := createUniqueAddressIndex(c.collection.account_balance); err != nil {
		return err
	}

	if err := createUniqueAddressIndex(c.collection.validator_stake); err != nil {
		return err
	}

	if err := createAddressTimeIndexDataIndex(c.collection.account_balance_index); err != nil {
		return err
	}

	if err := c.initMainnetGenesisBalance(); err != nil {
		return err
	}

	return nil
}

type addressInitInfo struct {
	address string `bson:"address"`
	balance int64  `bson:"balance"`
}

func (c *DbClient) initMainnetGenesisBalance() error {
	addressInit := []addressInitInfo{
		{"000000000000000000000000000000000000000000", 21000000000000000},
		{"pc1z2r0fmu8sg2ffa0tgrr08gnefcxl2kq7wvquf8z", 8400000000000000},
		{"pc1zprhnvcsy3pthekdcu28cw8muw4f432hkwgfasv", 6300000000000000},
		{"pc1znn2qxsugfrt7j4608zvtnxf8dnz8skrxguyf45", 4200000000000000},
		{"pc1zs64vdggjcshumjwzaskhfn0j9gfpkvche3kxd3", 2100000000000000},
	}

	for _, info := range addressInit {
		filter := bson.M{"address": info.address}
		var result bson.M
		err := c.collection.account_balance.FindOne(c.ctx, filter).Decode(&result)

		if err == mongo.ErrNoDocuments {
			doc := bson.M{
				"address": info.address,
				"balance": info.balance,
			}
			result, err := c.collection.account_balance.InsertOne(c.ctx, doc)
			if err != nil {
				return fmt.Errorf("failed to insert address %s: %v", info.address, err)
			}

			if result.InsertedID != nil {
				log.Printf("inserted address balance %s = %d", info.address, info.balance)
			}
		} else if err != nil {
			return fmt.Errorf("error checking address %s: %v", info.address, err)
		}
	}

	return nil
}

func (c *DbClient) updateAccountBalance(commitContext CommitContext) error {
	m := commitContext.GetTxMerger()

	if len(m.accountBalanceChange) == 0 {
		return nil
	}

	writeModel := make([]mongo.WriteModel, 0)

	for address, change := range m.accountBalanceChange {
		update := mongo.NewUpdateOneModel().
			SetFilter(bson.D{
				{Key: "address", Value: address},
			}).
			SetUpdate(bson.D{
				{Key: "$inc", Value: bson.D{
					{Key: "balance", Value: change},
				}},
			}).
			SetUpsert(true)

		writeModel = append(writeModel, update)
	}

	result, err := c.collection.account_balance.BulkWrite(c.ctx, writeModel)
	if err != nil {
		return err
	}

	if result.InsertedCount == 0 && result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		log.Printf("No data inserted or modified")
	}

	return nil
}

func (c *DbClient) updateAccountBalanceIndex(commitContext CommitContext) error {
	m := commitContext.GetTxMerger()

	if len(m.accountBalanceChange) == 0 {
		return nil
	}

	writeModel := make([]mongo.WriteModel, 0)

	for address, change := range m.accountBalanceChange {
		update := mongo.NewUpdateOneModel().
			SetFilter(bson.D{
				{Key: "time_index", Value: commitContext.GetTimeIndex()},
				{Key: "address", Value: address},
			}).
			SetUpdate(bson.D{
				{Key: "$inc", Value: bson.D{
					{Key: "balance", Value: change},
				}},
			}).
			SetUpsert(true)

		writeModel = append(writeModel, update)
	}

	result, err := c.collection.account_balance_index.BulkWrite(c.ctx, writeModel)
	if err != nil {
		return err
	}

	if result.InsertedCount == 0 && result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		log.Printf("No data inserted or modified")
	}

	return nil
}

func (c *DbClient) updateValidatorStake(commitContext CommitContext) error {
	m := commitContext.GetTxMerger()

	if len(m.validatorStakeChange) == 0 {
		return nil
	}

	writeModel := make([]mongo.WriteModel, 0)

	for address, change := range m.validatorStakeChange {
		var updateDoc bson.D

		if change > 0 {
			updateDoc = bson.D{
				{Key: "$inc", Value: bson.D{
					{Key: "stake", Value: change},
				}},
				{Key: "$inc", Value: bson.D{
					{Key: "stake_max", Value: change},
				}},
			}
		} else {
			updateDoc = bson.D{
				{Key: "$inc", Value: bson.D{
					{Key: "stake", Value: change},
				}},
			}
		}

		update := mongo.NewUpdateOneModel().
			SetFilter(bson.D{
				{Key: "address", Value: address},
			}).
			SetUpdate(updateDoc).
			SetUpsert(true)

		writeModel = append(writeModel, update)
	}

	result, err := c.collection.validator_stake.BulkWrite(c.ctx, writeModel)
	if err != nil {
		return err
	}

	if result.InsertedCount == 0 && result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		log.Printf("No data inserted or modified")
	}

	return nil
}
