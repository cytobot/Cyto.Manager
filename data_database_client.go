package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DatabaseClient struct {
	database          *gorm.DB
	CommandRepository *CommandRepository
}

func NewDatabaseClient(connectionString string) (*DatabaseClient, error) {
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	//defer db.Close()

	commandRepository, err := NewCommandRepository(db)
	if err != nil {
		return nil, err
	}

	client := &DatabaseClient{
		CommandRepository: commandRepository,
	}

	return client, nil
}
