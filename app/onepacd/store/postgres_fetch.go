package store

import (
	"context"

	"github.com/frimin/1pactus-react/app/onepacd/store/model"
)

func (s *postgresStore) GetNetworkGlobalStats(count int64) ([]model.GlobalState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), POSTGRES_DB_TIMEOUT)
	defer cancel()

	var rets []model.GlobalState

	if err := s.db.GetDB().WithContext(ctx).Order("time_index DESC").Limit(int(count)).Find(&rets).Error; err != nil {
		return nil, err
	}

	return rets, nil
}
