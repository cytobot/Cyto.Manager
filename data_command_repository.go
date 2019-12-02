package main

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type CommandRepository struct {
	db *gorm.DB
}

func NewCommandRepository(db *gorm.DB) (*CommandRepository, error) {
	if err := db.AutoMigrate(&CommandDefinition{}, &CommandParameterDefinition{}).Error; err != nil {
		return nil, err
	}

	repo := &CommandRepository{
		db: db,
	}

	if err := repo.seedInitialCommandDefinitions(); err != nil {
		return nil, err
	}

	return repo, nil
}

type CommandDefinition struct {
	CommandID            string `gorm:"primary_key"`
	Description          string
	Enabled              bool
	Unlisted             bool
	Triggers             pq.StringArray `gorm:"type:text[]"`
	PermissionLevel      string
	ParameterDefinitions []CommandParameterDefinition `gorm:"foreignKey:CommandID"`
	LastModifiedDateUtc  time.Time
	LastModifiedUserID   string
}

type CommandParameterDefinition struct {
	CommandID string `gorm:"primary_key"`
	Name      string `gorm:"primary_key"`
	Pattern   string
	Optional  bool
}

func (r *CommandRepository) Get(commandID string) (*CommandDefinition, error) {
	var commandDefinition *CommandDefinition

	if err := r.db.First(&commandDefinition).Error; err != nil {
		return nil, err
	}

	return commandDefinition, nil
}

func (r *CommandRepository) GetAll() ([]*CommandDefinition, error) {
	var commandDefinitions []*CommandDefinition

	if err := r.db.Find(&commandDefinitions).Error; err != nil {
		return nil, err
	}

	return commandDefinitions, nil
}

func (r *CommandRepository) Update(commandDefinition *CommandDefinition) (*CommandDefinition, error) {
	if err := r.db.Save(commandDefinition).Error; err != nil {
		return nil, err
	}

	return commandDefinition, nil
}

func (r *CommandRepository) Create(commandDefinition *CommandDefinition) (*CommandDefinition, error) {
	if err := r.db.Create(commandDefinition).Error; err != nil {
		return nil, err
	}

	return commandDefinition, nil
}

func (r *CommandRepository) Delete(commandDefinition CommandDefinition) error {
	if err := r.db.Delete(&commandDefinition).Error; err != nil {
		return err
	}

	return nil
}
