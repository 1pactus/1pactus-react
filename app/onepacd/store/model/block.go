package model

type Block struct {
	TimeIndex int64 `gorm:"primaryKey;uniqueIndex:idx_time_index;not null"`
	Height    int64 `gorm:"uniqueIndex:idx_height;not null"`
}
