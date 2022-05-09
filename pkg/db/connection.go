package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Init() (*sql.DB, error) {
	/* user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbName := os.Getenv("DB_NAME") */
	//connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbName)

	// Init database connection
	db, err := sql.Open("mysql", "root:21032991@tcp(localhost:3306)/products")
	if err != nil {
		return nil, err
	}

	// Check database connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
