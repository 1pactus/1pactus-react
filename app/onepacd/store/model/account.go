package model

// record of account balance latest state
type AccountBalance struct {
	Address string `gorm:"primaryKey;uniqueIndex:idx_address;not null"`
	Balance int64  `gorm:"not null"`
}

// record of account balance changes over timeindex
type AccountBalanceTimeIndex struct {
	Address       string `gorm:"primaryKey;not null"`
	TimeIndex     int64  `gorm:"primaryKey;not null"`
	BalanceChange int64  `gorm:"not null"`
}

// record of validator latest state
type ValidatorState struct {
	Address  string `gorm:"primaryKey;uniqueIndex:idx_address;not null"`
	Stake    int64  `gorm:"not null"`
	StakeMax int64  `gorm:"not null"`
}
