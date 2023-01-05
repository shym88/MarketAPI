package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Connect() {
	var err error

	dsn := os.Getenv("CONNECTION_STRING")
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully connected to the database")
	}
}

func CreateDB() {
	var err error

	dsn := os.Getenv("CONNECTION_STRING_POSTGRES")
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully connected to the database")
	}

	// check if db exists
	stmt := fmt.Sprintf("SELECT * FROM pg_database WHERE datname = '%s';", "marketapidb")
	rs := Database.Raw(stmt)
	if rs.Error != nil {
		panic(rs.Error)
	}

	// if not create it
	var rec = make(map[string]interface{})
	if rs.Find(rec); len(rec) == 0 {
		stmt := fmt.Sprintf("CREATE DATABASE %s;", "marketapidb")
		if rs := Database.Exec(stmt); rs.Error != nil {
			panic(rs.Error)
		}

		// close db connection
		sql, err := Database.DB()
		defer func() {
			_ = sql.Close()
		}()
		if err != nil {
			panic(err)
		}
	}

}
