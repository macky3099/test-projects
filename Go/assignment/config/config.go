package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = ""
	dbName   = "test"
)

var db *gorm.DB

func DatabaseConnection() *gorm.DB {

	if db != nil {
		return db
	}

	//sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

	sqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	db, err := gorm.Open(mysql.Open(sqlInfo), &gorm.Config{})
	if err != nil {
		return nil
	}

	return db
}

func GetConfig() *gorm.DB {
	return DatabaseConnection()
}
