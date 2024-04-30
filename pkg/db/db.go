package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func InitDB(opt Option) (*DB, error) {
	switch opt.Driver {
	case Mysql:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			opt.User, opt.Password, opt.Host, opt.Port, opt.DbName)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		return &DB{db}, nil
	default:
		return nil, errors.New("unsupported driver")
	}
}
