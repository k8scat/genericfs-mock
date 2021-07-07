package models

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wanhuasong/genericfs/config"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("mysql", config.Cfg.DSN)
	if err != nil {
		return err
	}
	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
	return nil
}

func Transact(do func(*sql.Tx) error) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	err = do(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
