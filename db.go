package common

import (
	"database/sql"
	"fmt"
	"third/gorm"
)

func InitDbPool(config *MysqlConfig) (*sql.DB, error) {

	dbPool, err := sql.Open("mysql", config.MysqlConn)
	if nil != err {
		Logger.Error("sql.Open error :%s,%v", config.MysqlConn, err)
		return nil, err
	}
	dbPool.SetMaxOpenConns(config.MysqlConnectPoolSize)
	dbPool.SetMaxIdleConns(config.MysqlConnectPoolSize)

	err = dbPool.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("init db pool OK")
	return dbPool, nil
}

func InitGormDbPool(config *MysqlConfig, setLog bool) (*gorm.DB, error) {

	db, err := gorm.Open("mysql", config.MysqlConn)
	if err != nil {
		fmt.Println("init db err : ", config, err)
		return nil, err
	}

	db.DB().SetMaxOpenConns(config.MysqlConnectPoolSize)
	db.DB().SetMaxIdleConns(config.MysqlConnectPoolSize)
	if setLog {
		db.LogMode(true)
		db.SetLogger(Logger)
	}
	db.SingularTable(true)

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}
	//	fmt.Println("init db pool OK")

	return &db, nil
}
