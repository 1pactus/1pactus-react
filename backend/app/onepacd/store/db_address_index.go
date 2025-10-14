package store

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createAddressDataIndex(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "address", Value: 1}},
		Options: options.Index().
			SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err
	}

	return nil
}

func createAddressTimeIndexDataIndex(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "time_index", Value: 1},
			{Key: "address", Value: 1},
		},
		Options: options.Index().
			SetUnique(false),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err
	}

	return nil
}

func (c *DbClient) createAddressIndex() error {
	c.collection.address_transfer_sender_index = c.database.Collection("address_transfer_sender_index")
	c.collection.address_transfer_receiver_index = c.database.Collection("address_transfer_receiver_index")
	c.collection.address_transfer_reward_index = c.database.Collection("address_transfer_reward_index")
	c.collection.address_bond_sender_index = c.database.Collection("address_bond_sender_index")
	c.collection.address_bond_receiver_index = c.database.Collection("address_bond_receiver_index")
	c.collection.address_unbond_index = c.database.Collection("address_unbond_index")
	c.collection.address_withdraw_sender_index = c.database.Collection("address_withdraw_sender_index")
	c.collection.address_withdraw_receiver_index = c.database.Collection("address_withdraw_receiver_index")

	if err := createAddressTimeIndexDataIndex(c.collection.address_transfer_sender_index); err != nil {
		return err
	}

	if err := createAddressTimeIndexDataIndex(c.collection.address_transfer_receiver_index); err != nil {
		return err
	}

	if err := createAddressTimeIndexDataIndex(c.collection.address_transfer_reward_index); err != nil {
		return err
	}

	if err := createAddressTimeIndexDataIndex(c.collection.address_bond_sender_index); err != nil {
		return err
	}

	if err := createAddressTimeIndexDataIndex(c.collection.address_bond_receiver_index); err != nil {
		return err
	}

	if err := createAddressTimeIndexDataIndex(c.collection.address_unbond_index); err != nil {
		return err
	}

	if err := createAddressTimeIndexDataIndex(c.collection.address_withdraw_sender_index); err != nil {
		return err
	}

	if err := createAddressTimeIndexDataIndex(c.collection.address_withdraw_receiver_index); err != nil {
		return err
	}

	return nil
}

func (c *DbClient) updateSenderTransfers(commitContext CommitContext) error {
	m := commitContext.GetTxMerger()

	if len(m.transferSender) == 0 {
		return nil
	}

	writeModel := make([]mongo.WriteModel, 0)

	for timeIndex, record := range m.transferSender {
		for sender, transferSender := range record {
			incFields := bson.D{}

			incFields = append(incFields, bson.E{Key: "sender_amount", Value: transferSender.total})

			for targetAddress, amount := range transferSender.addresses {
				key := fmt.Sprintf("sender_map.%s", targetAddress)
				incFields = append(incFields, bson.E{Key: key, Value: amount})
			}

			update := mongo.NewUpdateOneModel().
				SetFilter(bson.D{
					{Key: "time_index", Value: timeIndex},
					{Key: "address", Value: sender},
				}).
				SetUpdate(bson.D{
					{Key: "$inc", Value: incFields},
				}).
				SetUpsert(true)

			writeModel = append(writeModel, update)
		}
	}

	result, err := c.collection.address_transfer_sender_index.BulkWrite(c.ctx, writeModel)
	if err != nil {
		return err
	}

	if result.InsertedCount == 0 && result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		log.Printf("updateSenderTransfers: No data inserted or modified")
	}

	return nil
}

func (c *DbClient) updateReceiverTransfers(commitContext CommitContext) error {
	m := commitContext.GetTxMerger()

	if len(m.transferReceiver) == 0 {
		return nil
	}

	writeModel := make([]mongo.WriteModel, 0)

	for timeIndex, record := range m.transferReceiver {
		for receiver, transferReceiver := range record {
			incFields := bson.D{}

			incFields = append(incFields, bson.E{Key: "receiver_amount", Value: transferReceiver.total})

			for targetAddress, amount := range transferReceiver.addresses {
				key := fmt.Sprintf("receiver_map.%s", targetAddress)
				incFields = append(incFields, bson.E{Key: key, Value: amount})
			}

			update := mongo.NewUpdateOneModel().
				SetFilter(bson.D{
					{Key: "time_index", Value: timeIndex},
					{Key: "address", Value: receiver},
				}).
				SetUpdate(bson.D{
					{Key: "$inc", Value: incFields},
				}).
				SetUpsert(true)

			writeModel = append(writeModel, update)
		}
	}

	result, err := c.collection.address_transfer_receiver_index.BulkWrite(c.ctx, writeModel)
	if err != nil {
		return err
	}

	if result.InsertedCount == 0 && result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		log.Printf("updateReceiverTransfers: No data inserted or modified")
	}

	return nil
}

func (c *DbClient) updateRewardTransfers(commitContext CommitContext) error {
	m := commitContext.GetTxMerger()

	if len(m.transferReward) == 0 {
		return nil
	}

	writeModel := make([]mongo.WriteModel, 0)

	for timeIndex, record := range m.transferReward {
		for rewardAddress, transferReward := range record {

			// 处理 sender_map
			incFields := bson.D{}

			incFields = append(incFields, bson.E{Key: "reward_amount", Value: transferReward.total})

			for targetAddress, amount := range transferReward.addresses {
				key := fmt.Sprintf("reward_map.%s", targetAddress)
				incFields = append(incFields, bson.E{Key: key, Value: amount})
			}

			update := mongo.NewUpdateOneModel().
				SetFilter(bson.D{
					{Key: "time_index", Value: timeIndex},
					{Key: "address", Value: rewardAddress},
				}).
				SetUpdate(bson.D{
					{Key: "$inc", Value: incFields},
				}).
				SetUpsert(true)

			writeModel = append(writeModel, update)
		}
	}

	// Execute bulk write
	result, err := c.collection.address_transfer_reward_index.BulkWrite(c.ctx, writeModel)
	if err != nil {
		return err
	}

	if result.InsertedCount == 0 && result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		log.Printf("updateRewardTransfers: No data inserted or modified")
	}

	return nil
}

func (c *DbClient) updateSenderBond(commitContext CommitContext) error {
	m := commitContext.GetTxMerger()

	if len(m.bondSender) == 0 {
		return nil
	}

	writeModel := make([]mongo.WriteModel, 0)

	for timeIndex, record := range m.bondSender {
		for sender, bondSender := range record {
			incFields := bson.D{}

			incFields = append(incFields, bson.E{Key: "sender_amount", Value: bondSender.total})

			for targetAddress, amount := range bondSender.addresses {
				key := fmt.Sprintf("sender_map.%s", targetAddress)
				incFields = append(incFields, bson.E{Key: key, Value: amount})
			}

			update := mongo.NewUpdateOneModel().
				SetFilter(bson.D{
					{Key: "time_index", Value: timeIndex},
					{Key: "address", Value: sender},
				}).
				SetUpdate(bson.D{
					{Key: "$inc", Value: incFields},
				}).
				SetUpsert(true)

			writeModel = append(writeModel, update)
		}
	}

	// Execute bulk write
	result, err := c.collection.address_bond_sender_index.BulkWrite(c.ctx, writeModel)
	if err != nil {
		return err
	}

	if result.InsertedCount == 0 && result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		log.Printf("updateSenderBond: No data inserted or modified")
	}

	return nil
}

func (c *DbClient) updateReceiverBond(commitContext CommitContext) error {
	m := commitContext.GetTxMerger()

	if len(m.bondReceiver) == 0 {
		return nil
	}

	writeModel := make([]mongo.WriteModel, 0)

	for timeIndex, record := range m.bondReceiver {
		for receiver, bondReceiver := range record {
			incFields := bson.D{}

			incFields = append(incFields, bson.E{Key: "receiver_amount", Value: bondReceiver.total})

			for targetAddress, amount := range bondReceiver.addresses {
				key := fmt.Sprintf("receiver_map.%s", targetAddress)
				incFields = append(incFields, bson.E{Key: key, Value: amount})
			}

			update := mongo.NewUpdateOneModel().
				SetFilter(bson.D{
					{Key: "time_index", Value: timeIndex},
					{Key: "address", Value: receiver},
				}).
				SetUpdate(bson.D{
					{Key: "$inc", Value: incFields},
				}).
				SetUpsert(true)

			writeModel = append(writeModel, update)
		}
	}

	// Execute bulk write
	result, err := c.collection.address_bond_receiver_index.BulkWrite(c.ctx, writeModel)
	if err != nil {
		return err
	}

	if result.InsertedCount == 0 && result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		log.Printf("updateReceiverBond: No data inserted or modified")
	}

	return nil
}

func (c *DbClient) updateUnbondTransfers(commitContext CommitContext) error {
	m := commitContext.GetTxMerger()

	if len(m.unbond) == 0 {
		return nil
	}

	writeModel := make([]mongo.WriteModel, 0)
	writeValidatorModel := make([]mongo.WriteModel, 0)

	for timeIndex, record := range m.unbond {
		for rewardAddress, info := range record {
			insert := mongo.NewInsertOneModel().
				SetDocument(bson.D{
					{Key: "time_index", Value: timeIndex},
					{Key: "address", Value: rewardAddress},
					{Key: "height", Value: info.Height},
					{Key: "hash", Value: info.Hash},
					{Key: "time", Value: info.Time},
				})

			writeModel = append(writeModel, insert)

			validatorInsert := mongo.NewUpdateOneModel().
				SetFilter(bson.D{
					{Key: "address", Value: rewardAddress},
				}).
				SetUpdate(bson.D{
					{Key: "$set", Value: bson.M{
						"unbond_time_index": timeIndex,
						"unbond_height":     info.Height,
						"unbond_hash":       info.Hash,
						"unbond_time":       info.Time,
					}},
				}).
				SetUpsert(true)

			writeValidatorModel = append(writeValidatorModel, validatorInsert)
		}
	}

	// Execute bulk write
	result, err := c.collection.address_unbond_index.BulkWrite(c.ctx, writeModel)
	if err != nil {
		return err
	}

	if result.InsertedCount == 0 && result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		log.Printf("updateRewardTransfers: No data inserted or modified")
	}

	result1, err := c.collection.validator_stake.BulkWrite(c.ctx, writeValidatorModel)
	if err != nil {
		return err
	}

	if result1.InsertedCount == 0 && result1.ModifiedCount == 0 && result1.UpsertedCount == 0 {
		log.Printf("updateUnbondTransfers.writeValidatorModel: No data inserted or modified")
	}

	return nil
}

func (c *DbClient) updateWithdrawSender(commitContext CommitContext) error {
	m := commitContext.GetTxMerger()

	if len(m.withdrawSender) == 0 {
		return nil
	}

	writeModel := make([]mongo.WriteModel, 0)

	for timeIndex, record := range m.withdrawSender {
		for validator, withdrawSender := range record {
			incFields := bson.D{}

			incFields = append(incFields, bson.E{Key: "sender_amount", Value: withdrawSender.total})

			for targetAddress, amount := range withdrawSender.addresses {
				key := fmt.Sprintf("sender_map.%s", targetAddress)
				incFields = append(incFields, bson.E{Key: key, Value: amount})
			}

			update := mongo.NewUpdateOneModel().
				SetFilter(bson.D{
					{Key: "time_index", Value: timeIndex},
					{Key: "address", Value: validator},
				}).
				SetUpdate(bson.D{
					{Key: "$inc", Value: incFields},
				}).
				SetUpsert(true)

			writeModel = append(writeModel, update)
		}
	}

	result, err := c.collection.address_withdraw_sender_index.BulkWrite(c.ctx, writeModel)
	if err != nil {
		return err
	}

	if result.InsertedCount == 0 && result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		log.Printf("updateWithdrawSender: No data inserted or modified")
	}

	return nil
}

func (c *DbClient) updateWithdrawReceiver(commitContext CommitContext) error {
	m := commitContext.GetTxMerger()

	if len(m.withdrawReceiver) == 0 {
		return nil
	}

	writeModel := make([]mongo.WriteModel, 0)

	for timeIndex, record := range m.withdrawReceiver {
		for account, withdrawReceiver := range record {
			incFields := bson.D{}

			incFields = append(incFields, bson.E{Key: "receiver_amount", Value: withdrawReceiver.total})

			for targetAddress, amount := range withdrawReceiver.addresses {
				key := fmt.Sprintf("receiver_map.%s", targetAddress)
				incFields = append(incFields, bson.E{Key: key, Value: amount})
			}

			update := mongo.NewUpdateOneModel().
				SetFilter(bson.D{
					{Key: "time_index", Value: timeIndex},
					{Key: "address", Value: account},
				}).
				SetUpdate(bson.D{
					{Key: "$inc", Value: incFields},
				}).
				SetUpsert(true)

			writeModel = append(writeModel, update)
		}
	}

	result, err := c.collection.address_withdraw_receiver_index.BulkWrite(c.ctx, writeModel)
	if err != nil {
		return err
	}

	if result.InsertedCount == 0 && result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		log.Printf("updateWithdrawReceiver: No data inserted or modified")
	}

	return nil
}
