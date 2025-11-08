package storedriver

import (
	"context"
	"fmt"
	"time"

	"github.com/1pactus/1pactus-react/config"
	"github.com/1pactus/1pactus-react/log"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type IPostgresGormStore interface {
	Init(store GormPostgres)
	AutoMigrate() error
	Models() []interface{}
	Indexes() []IndexSchema
}

type TableSchema struct {
	Name    string
	DDL     string
	Indexes []IndexSchema
}

type IndexSchema struct {
	Name    string
	Table   string
	Columns []string
	Unique  bool
}

type GormPostgres interface {
	GetDB() *gorm.DB
	GetTimeout() time.Duration
	WithContext(ctx context.Context) *gorm.DB
}

type postgresGormImpl struct {
	conf    *config.PostgresConfig
	db      *gorm.DB
	stores  []IPostgresGormStore
	timeout time.Duration
	log     log.ILogger
}

func (db *postgresGormImpl) Close() {
	if db.db != nil {
		db, err := db.db.DB()
		if err == nil {
			db.Close()
		}
	}
}

func PostgresGormStart(name string, conf *config.PostgresConfig, stores []IPostgresGormStore) error {
	m := &postgresGormImpl{
		conf:    conf,
		stores:  stores,
		timeout: time.Second * 10, // Default timeout
		log:     log.WithKv("module", "store").WithKv("postgres", name),
	}

	var err error
	maxRetry := 10

	for {
		if err = m.connect(); err == nil {
			m.initStores()

			m.log.Infof("gorm postgres connect and initialized success")

			if err := m.autoMigrate(); err != nil {
				return fmt.Errorf("gorm postgres [%s] auto migrate failed: %v", name, err)
			}

			if err := m.ensureIndexes(); err != nil {
				return fmt.Errorf("gorm postgres [%s] ensure indexes failed: %v", name, err)
			}

			go m.monitorConnection()
			return nil
		} else {
			maxRetry -= 1

			if maxRetry <= 0 {
				return fmt.Errorf("gorm postgres [%s] connect failed after retries: %v", name, err)
			}

			m.log.Errorf("gorm postgres [%s] connect failed: %v", name, err)
			time.Sleep(time.Second * 5)
		}
	}
}

func (db *postgresGormImpl) connect() error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		db.conf.Host,
		db.conf.Port,
		db.conf.Username,
		db.conf.Password,
		db.conf.Database)

	var gormLogLevel logger.LogLevel
	switch db.log.GetInternalLogger().GetLevel() {
	case zerolog.DebugLevel:
		gormLogLevel = logger.Info
	case zerolog.ErrorLevel:
		gormLogLevel = logger.Error
	default:
		gormLogLevel = logger.Warn
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return err
	}

	postgresDb, err := gormDB.DB()
	if err != nil {
		return err
	}

	postgresDb.SetMaxOpenConns(db.conf.MaxOpenConns)
	postgresDb.SetMaxIdleConns(db.conf.MaxIdleConns)
	postgresDb.SetConnMaxLifetime(time.Duration(db.conf.ConnMaxLifetime) * time.Second)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := postgresDb.PingContext(ctx); err != nil {
		postgresDb.Close()
		return err
	}

	db.db = gormDB
	return nil
}

func (db *postgresGormImpl) initStores() {
	for _, store := range db.stores {
		store.Init(db)
	}
}

func (db *postgresGormImpl) autoMigrate() error {
	for _, store := range db.stores {
		// 使用GORM的AutoMigrate功能
		if err := store.AutoMigrate(); err != nil {
			return fmt.Errorf("auto migrate failed for store: %v", err)
		}

		// 如果store有自定义的Models方法，也可以使用GORM的AutoMigrate
		models := store.Models()
		if len(models) > 0 {
			if err := db.db.AutoMigrate(models...); err != nil {
				return fmt.Errorf("auto migrate models failed: %v", err)
			}
			db.log.Infof("auto migrated %d models successfully", len(models))
		}
	}
	return nil
}

func (db *postgresGormImpl) ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	for _, store := range db.stores {
		indexes := store.Indexes()
		for _, index := range indexes {
			if err := db.ensureIndex(ctx, index); err != nil {
				return err
			}
		}
	}
	return nil
}

func (db *postgresGormImpl) ensureIndex(ctx context.Context, index IndexSchema) error {
	migrator := db.db.Migrator()

	if migrator.HasIndex(index.Table, index.Name) {
		db.log.Infof("index %s already exists", index.Name)
		return nil
	}

	var indexType string
	if index.Unique {
		indexType = "UNIQUE INDEX"
	} else {
		indexType = "INDEX"
	}

	columns := ""
	for i, col := range index.Columns {
		if i > 0 {
			columns += ", "
		}
		columns += fmt.Sprintf("`%s`", col)
	}

	createIndexSQL := fmt.Sprintf("CREATE %s `%s` ON `%s` (%s)",
		indexType, index.Name, index.Table, columns)

	if err := db.db.WithContext(ctx).Exec(createIndexSQL).Error; err != nil {
		return fmt.Errorf("create index %s failed: %v", index.Name, err)
	}

	db.log.Infof("index %s created successfully", index.Name)
	return nil
}

func (db *postgresGormImpl) monitorConnection() {
	healthcheck := time.Duration(db.conf.Healthcheck)

	if healthcheck <= 1 {
		healthcheck = 1
	}

	for {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			sqlDB, err := db.db.DB()
			if err != nil {
				db.log.Errorf("failed to get underlying sql.DB: %v", err)
				return
			}

			if err := sqlDB.PingContext(ctx); err != nil {
				db.log.Errorf("lost pgsql connection, retrying: %v", err)

				for {
					if err := db.connect(); err == nil {
						db.initStores()
						db.log.Info("successfully reconnected to pgsql")
						break
					}
					db.log.Errorf("error connecting to pgsql: %v", err)
					time.Sleep(5 * time.Second)
				}
			}
		}()

		time.Sleep(healthcheck * time.Second)
	}
}

func (db *postgresGormImpl) GetDB() *gorm.DB {
	return db.db
}

func (db *postgresGormImpl) GetTimeout() time.Duration {
	return db.timeout
}

func (db *postgresGormImpl) WithContext(ctx context.Context) *gorm.DB {
	return db.db.WithContext(ctx)
}
