package store

import (
	_ "embed"
	"time"

	"github.com/frimin/1pactus-react/app/onepacd/store/model"
	"github.com/frimin/1pactus-react/store/storedriver"
)

const (
	POSTGRES_DB_TIMEOUT = 10 * time.Second
)

type postgresStore struct {
	db storedriver.GormPostgres
}

func (s *postgresStore) Init(db storedriver.GormPostgres) {
	s.db = db
}

func (s *postgresStore) AutoMigrate() error {
	return s.db.GetDB().AutoMigrate(
		&model.GlobalState{},
		&model.Block{},
	)
}

func (s *postgresStore) Models() []interface{} {
	return []interface{}{
		&model.GlobalState{},
		&model.Block{},
	}
}

func (s *postgresStore) Indexes() []storedriver.IndexSchema {
	return []storedriver.IndexSchema{}
}
