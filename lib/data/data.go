package data

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DatabaseClient struct {
	database          *gorm.DB
	CommandRepository *CommandRepository
}

func NewDatabaseClient(dialect string, connectionString string) (*DatabaseClient, error) {
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	commandRepository, err := NewCommandRepository(db)
	if err != nil {
		return nil, err
	}

	client := &DatabaseClient{
		CommandRepository: commandRepository,
	}

	return client, nil
}
