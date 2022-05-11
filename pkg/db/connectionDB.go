package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Init() (*sql.DB, error) {
	os.Setenv("USER", "root")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbName)
	fmt.Println(connectionString)
	// Init database connection
	// db, err := sql.Open("mysql", "root:21032991@tcp(localhost:3306)/products")
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Check database connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connection database success")

	return db, nil
}
