package model

import (
	"fmt"
)

const (
	TreasuryAddress = "000000000000000000000000000000000000000000"
)

type TxUnbond struct {
	Height int64
	Time   int64
	Hash   string
}

type TxMerger struct {
	transferReceiver map[int64]map[string]*txTransferMerged
	transferSender   map[int64]map[string]*txTransferMerged
	transferReward   map[int64]map[string]*txTransferMerged

	bondReceiver map[int64]map[string]*txTransferMerged
	bondSender   map[int64]map[string]*txTransferMerged

	unbond map[int64]map[string]TxUnbond

	withdrawSender   map[int64]map[string]*txTransferMerged
	withdrawReceiver map[int64]map[string]*txTransferMerged

	accountBalanceChange map[string]int64
	validatorStakeChange map[string]int64
}

type txTransferMerged struct {
	total     int64
	addresses map[string]int64
}

func NewTxMerger() *TxMerger {
	m := &TxMerger{}
	m.Clean()

	return m
}

func (m *TxMerger) AddTransfer(timeIndex int64, sender string, receiver string, amount int64, fee int64) error {
	if _, ok := m.transferReceiver[timeIndex]; !ok {
		m.transferReceiver[timeIndex] = make(map[string]*txTransferMerged)
	}

	if _, ok := m.transferSender[timeIndex]; !ok {
		m.transferSender[timeIndex] = make(map[string]*txTransferMerged)
	}

	if record, ok := m.transferReceiver[timeIndex][receiver]; !ok {
		m.transferReceiver[timeIndex][receiver] = &txTransferMerged{
			total:     amount,
			addresses: map[string]int64{sender: amount},
		}
	} else {
		record.total += amount
		record.addresses[sender] += amount
	}

	if record, ok := m.transferSender[timeIndex][sender]; !ok {
		m.transferSender[timeIndex][sender] = &txTransferMerged{
			total:     amount,
			addresses: map[string]int64{receiver: amount},
		}
	} else {
		record.total += amount
		record.addresses[receiver] += amount
	}

	if _, ok := m.accountBalanceChange[sender]; !ok {
		m.accountBalanceChange[sender] = -(amount + fee)
	} else {
		m.accountBalanceChange[sender] -= (amount + fee)
	}

	if _, ok := m.accountBalanceChange[receiver]; !ok {
		m.accountBalanceChange[receiver] = amount
	} else {
		m.accountBalanceChange[receiver] += amount
	}

	return nil
}

func (m *TxMerger) AddReward(timeIndex int64, receiver string, amount int64, proposerAddress string) error {
	if _, ok := m.transferReward[timeIndex]; !ok {
		m.transferReward[timeIndex] = make(map[string]*txTransferMerged)
	}

	if record, ok := m.transferReward[timeIndex][receiver]; !ok {
		m.transferReward[timeIndex][receiver] = &txTransferMerged{
			total:     amount,
			addresses: map[string]int64{proposerAddress: amount},
		}
	} else {
		record.total += amount
		record.addresses[proposerAddress] += amount
	}

	if _, ok := m.accountBalanceChange[receiver]; !ok {
		m.accountBalanceChange[receiver] = amount
	} else {
		m.accountBalanceChange[receiver] += amount
	}

	if _, ok := m.accountBalanceChange[receiver]; !ok {
		m.accountBalanceChange[TreasuryAddress] = -amount
	} else {
		m.accountBalanceChange[TreasuryAddress] -= amount
	}

	return nil
}

func (m *TxMerger) AddBond(timeIndex int64, sender string, receiver string, stake int64, fee int64) error {
	if _, ok := m.bondReceiver[timeIndex]; !ok {
		m.bondReceiver[timeIndex] = make(map[string]*txTransferMerged)
	}

	if _, ok := m.bondSender[timeIndex]; !ok {
		m.bondSender[timeIndex] = make(map[string]*txTransferMerged)
	}

	if record, ok := m.bondReceiver[timeIndex][receiver]; !ok {
		m.bondReceiver[timeIndex][receiver] = &txTransferMerged{
			total:     stake,
			addresses: map[string]int64{sender: stake},
		}
	} else {
		record.total += stake
		record.addresses[sender] += stake
	}

	if record, ok := m.bondSender[timeIndex][sender]; !ok {
		m.bondSender[timeIndex][sender] = &txTransferMerged{
			total:     stake,
			addresses: map[string]int64{receiver: stake},
		}
	} else {
		record.total += stake
		record.addresses[receiver] += stake
	}

	if _, ok := m.accountBalanceChange[sender]; !ok {
		m.accountBalanceChange[sender] = -(stake + fee)
	} else {
		m.accountBalanceChange[sender] -= (stake + fee)
	}

	if _, ok := m.validatorStakeChange[receiver]; !ok {
		m.validatorStakeChange[receiver] = stake
	} else {
		m.validatorStakeChange[receiver] += stake
	}

	return nil
}

func (m *TxMerger) AddUnbond(timeIndex int64, validator string, height int64, hash string, time int64) error {
	if _, ok := m.unbond[timeIndex]; !ok {
		m.unbond[timeIndex] = make(map[string]TxUnbond)
	}

	if _, ok := m.unbond[timeIndex][validator]; !ok {
		m.unbond[timeIndex][validator] = TxUnbond{Height: height, Hash: hash, Time: time}
	} else {
		return fmt.Errorf("unbond record already exists")
	}

	return nil
}

func (m *TxMerger) AddWithdraw(timeIndex int64, validator string, account string, amount int64, fee int64) error {
	if _, ok := m.withdrawSender[timeIndex]; !ok {
		m.withdrawSender[timeIndex] = make(map[string]*txTransferMerged)
	}

	if _, ok := m.withdrawReceiver[timeIndex]; !ok {
		m.withdrawReceiver[timeIndex] = make(map[string]*txTransferMerged)
	}

	if record, ok := m.withdrawSender[timeIndex][validator]; !ok {
		m.withdrawSender[timeIndex][validator] = &txTransferMerged{
			total:     amount,
			addresses: map[string]int64{account: amount},
		}
	} else {
		record.total += amount
		record.addresses[account] += amount
	}

	if record, ok := m.withdrawReceiver[timeIndex][account]; !ok {
		m.withdrawReceiver[timeIndex][account] = &txTransferMerged{
			total:     amount,
			addresses: map[string]int64{validator: amount},
		}
	} else {
		record.total += amount
		record.addresses[validator] += amount
	}

	if _, ok := m.accountBalanceChange[account]; !ok {
		m.accountBalanceChange[account] = amount
	} else {
		m.accountBalanceChange[account] += amount
	}

	if _, ok := m.validatorStakeChange[validator]; !ok {
		m.validatorStakeChange[validator] = -(amount + fee)
	} else {
		m.validatorStakeChange[validator] -= (amount + fee)
	}

	return nil
}

func (m *TxMerger) Clean() {
	m.transferReceiver = make(map[int64]map[string]*txTransferMerged)
	m.transferSender = make(map[int64]map[string]*txTransferMerged)
	m.transferReward = make(map[int64]map[string]*txTransferMerged)
	m.bondReceiver = make(map[int64]map[string]*txTransferMerged)
	m.bondSender = make(map[int64]map[string]*txTransferMerged)
	m.withdrawReceiver = make(map[int64]map[string]*txTransferMerged)
	m.withdrawSender = make(map[int64]map[string]*txTransferMerged)
	m.unbond = make(map[int64]map[string]TxUnbond)
	m.accountBalanceChange = make(map[string]int64)
	m.validatorStakeChange = make(map[string]int64)
}
