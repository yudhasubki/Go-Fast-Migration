package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yudhasubki/go-fastmigration/migrations"
)

func main() {
	db, err := sql.Open("mysql", "root:test123@tcp(127.0.0.1:3306)/test_database")
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)
	migrations.MigrateContainer(db)
}
