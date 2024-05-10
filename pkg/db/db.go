package db

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type DriverType string

const (
	Postgres DriverType = "postgres"
)

type Option struct {
	Driver          DriverType `json:"driver" yaml:"driver"`
	Host            string     `json:"host" yaml:"host"`
	Port            string     `json:"port" yaml:"port"`
	User            string     `json:"user" yaml:"user"`
	Password        string     `json:"password" yaml:"password"`
	DbName          string     `json:"db_name" yaml:"db_name"`
	MaxOpenConns    int        `json:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int        `json:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifeTime int        `json:"conn_max_life_time" yaml:"conn_max_life_time"`
	Debug           bool
}

type DB struct {
	*gorm.DB
	sqlDB *sql.DB
}

func NewDB(opt Option) (*DB, error) {
	var (
		db  *gorm.DB
		err error
	)
	switch opt.Driver {
	case Postgres:
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			opt.Host, opt.User, opt.Password, opt.DbName, opt.Port)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, errors.New("unsupported driver")
	}
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	maxOpenConns := opt.MaxOpenConns
	if maxOpenConns == 0 {
		maxOpenConns = 10
	}
	maxIdleConns := opt.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = 5
	}
	connMaxLifeTime := opt.ConnMaxLifeTime
	if opt.ConnMaxLifeTime == 0 {
		connMaxLifeTime = 60000
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifeTime) * time.Millisecond)
	if opt.Debug {
		db = db.Debug()
	}
	return &DB{DB: db, sqlDB: sqlDB}, nil
}

func (db *DB) Close() error {
	return db.sqlDB.Close()
}
