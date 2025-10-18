package store

import (
	"errors"

	"github.com/frimin/1pactus-react/app/onepacd/store/model"
	"gorm.io/gorm"
)

func (s *postgresStore) GetTopGlobalState() (*model.GlobalState, error) {
	state := model.NewGlobalState()
	err := s.db.GetDB().Order("time_index desc").First(state).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return state, nil
}

func (s *postgresStore) InsertGlobalState(state *model.GlobalState) error {
	return s.db.GetDB().Create(state).Error
}

func (s *postgresStore) GetTopBlock() (*model.Block, error) {
	block := &model.Block{}
	err := s.db.GetDB().Order("height desc").First(block).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return block, nil
}
