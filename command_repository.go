package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type CommandRepository struct {
	db *gorm.DB
}

func NewCommandRepository(db *gorm.DB) (*CommandRepository, error) {
	if err := db.AutoMigrate(&CommandDefinition{}).Error; err != nil {
		return nil, err
	}

	return &CommandRepository{
		db: db,
	}, nil
}

type CommandDefinition struct {
	gorm.Model
	CommandID            string `gorm:"primary_key"`
	Enabled              bool
	Triggers             []string `gorm:"type:text[]"`
	PermissionLevel      string
	ParameterDefinitions []CommandParameterDefinition `gorm:"foreignKey:CommandID"`
	LastModifiedDateUtc  time.Time
	LastModifiedUserID   string
}

type CommandParameterDefinition struct {
	gorm.Model
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

func (r *CommandRepository) Create(commandDefinition CommandDefinition) (*CommandDefinition, error) {
	if err := r.db.Create(&commandDefinition).Error; err != nil {
		return nil, err
	}

	return &commandDefinition, nil
}

func (r *CommandRepository) Delete(commandDefinition CommandDefinition) error {
	if err := r.db.Delete(&commandDefinition).Error; err != nil {
		return err
	}

	return nil
}
