package database

import (
	"fmt"
	"net/url"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

// Database instance, implements IDatabase
type Database struct {
	db     *gorm.DB
	logger *zap.Logger
	config *DBConfig
}

func NewDatabaseClient(logger *zap.Logger, conf *DBConfig, skipMigration bool) (*Database, error) {
	// create db connection
	dsn := url.URL{
		User:     url.UserPassword(conf.User, conf.Password),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%v", conf.Host, conf.Port),
		Path:     conf.DBName,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}

	level := logger2.Silent
	if logger.Core().Enabled(zap.DebugLevel) {
		level = logger2.Warn
	}

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{Logger: zapgorm2.New(logger).LogMode(level)})
	if err != nil {
		logger.Debug("failed to connect to db", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to db, error: %w", err)
	}

	if !skipMigration {
		// initialize tables
		if err := db.AutoMigrate(&MintRequest{}); err != nil {
			return nil, fmt.Errorf("failed to migrate mint request table, error: %w", err)
		}
	}
	return &Database{
		db:     db,
		logger: logger,
		config: conf,
	}, nil
}

func (d *Database) Insert(tx *gorm.DB, payload *MintRequest) error {
	var res *gorm.DB
	if tx == nil {
		res = d.db.Create(payload)
	} else {
		res = tx.Create(payload)
	}
	return res.Error
}

func (d *Database) Find(tx *gorm.DB, payload *MintRequest) ([]MintRequest, error) {
	requests := []MintRequest{}
	var res *gorm.DB
	if tx == nil {
		res = d.db.Where(payload).Find(&requests)
	} else {
		res = tx.Where(payload).Find(&requests)
	}
	return requests, res.Error
}

func (d *Database) Update(tx *gorm.DB, payload *MintRequest, all bool) error {
	var res *gorm.DB
	if tx == nil {
		tempDB := d.db.Model(payload)
		if all {
			tempDB = tempDB.Select("*")
		}
		res = tempDB.Updates(payload)
	} else {
		tempDB := tx.Model(payload)
		if all {
			tempDB = tempDB.Select("*")
		}
		res = tempDB.Updates(payload)
	}
	return res.Error
}

func (d *Database) Delete(tx *gorm.DB, payload *MintRequest) error {
	var res *gorm.DB
	if tx == nil {
		res = d.db.Delete(payload)
	} else {
		res = tx.Delete(payload)
	}
	return res.Error
}

func (d *Database) Transaction(tx func(tx *gorm.DB) error) error {
	return d.db.Transaction(tx)
}

func (d *Database) DropTable() error {
	return d.db.Migrator().DropTable(&MintRequest{})
}
