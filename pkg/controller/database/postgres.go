package database

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type DB struct {
	*sql.DB
}

func open() (*DB, error) {
	dbConfig := viper.Sub("postgres")

	connStr := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable",
		dbConfig.GetString("database"),
		dbConfig.GetString("user"),
		dbConfig.GetString("password"),
		dbConfig.GetString("host"))

	log.Info(connStr)
	db, err := sql.Open("postgres", connStr)

	return &DB{db}, err
}

func postgresCreateDB(dbName string) error {
	db, err := open()
	if err != nil {
		log.Error(err, "Unable to connect to the database")
	}
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", dbName))
	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Code == "42P04" {
			log.Error(err, "Database already exists, wan't do anything")
			return nil
		}
		log.Error(err, "Unable to create database", "Database:", dbName)
	}
	return err
}

func postgresCreateUser(userName string) (string, error) {
	return "nil", nil
}

func postgresDelUser(userName string) (string, error) {
	return "nil", nil
}

func postgresDelDB(dbName string) error {
	db, err := open()
	if err != nil {
		log.Error(err, "Unable to connect to the database")
	}
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE \"%s\"", dbName))
	if err != nil {
		log.Error(err, "Unable to drop the database", "Database:", dbName)
	}

	return err
}

func grant(dbName string, userName string) error {
	return nil
}
