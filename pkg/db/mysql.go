package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysqlDB() *gorm.DB {
	os.Setenv("USER", "root")
	user := os.Getenv("USER")
	port := os.Getenv("PORT")
	charset := os.Getenv("CHARSET")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	dbName := os.Getenv("DB_NAME")
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", user, password, host, port, dbName, charset)
	db, err := gorm.Open(mysql.Open(connString), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	return db
}
