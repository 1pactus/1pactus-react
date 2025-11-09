package model

type TxTransferTimeIndex struct {
	AddressFrom string `gorm:"primaryKey;not null"`
	AddressTo   string `gorm:"primaryKey;not null"`
	TimeIndex   int64  `gorm:"primaryKey;not null"`
	Amount      int64  `gorm:"not null"`
}

type TxRewardTimeIndex struct {
	Address   string `gorm:"primaryKey;not null"`
	TimeIndex int64  `gorm:"primaryKey;not null"`
	Amount    int64  `gorm:"not null"`
}

type TxBondTimeIndex struct {
	AddressFrom string `gorm:"primaryKey;not null"`
	AddressTo   string `gorm:"primaryKey;not null"`
	TimeIndex   int64  `gorm:"primaryKey;not null"`
	Amount      int64  `gorm:"not null"`
}

type TxUnbondTimeIndex struct {
	Address   string `gorm:"primaryKey;not null"`
	TimeIndex int64  `gorm:"primaryKey;not null"`
	Time      int64  `gorm:"not null"`
	Hash      string `gorm:"not null"`
}

type TxWithdrawTimeIndex struct {
	AddressFrom string `gorm:"primaryKey;not null"`
	AddressTo   string `gorm:"primaryKey;not null"`
	TimeIndex   int64  `gorm:"primaryKey;not null"`
	Amount      int64  `gorm:"not null"`
}
