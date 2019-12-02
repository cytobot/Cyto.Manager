package main

import (
	"fmt"
)

func (s *ManagerState) getAllCommandDefinitions() ([]*CommandDefinition, error) {
	commandDefinitions, err := s.data.CommandRepository.GetAll()

	return commandDefinitions, err
}

func (s *ManagerState) getCommandDefinition(commandID string) (*CommandDefinition, error) {
	commandDefinition, err := s.data.CommandRepository.Get(commandID)

	return commandDefinition, err
}

func (s *ManagerState) createCommandDefinition(commandDefinition *CommandDefinition) (*CommandDefinition, error) {
	existingCommandDefinition, err := s.data.CommandRepository.Get(commandDefinition.CommandID)
	if err != nil {
		return nil, err
	}

	if existingCommandDefinition != nil {
		return nil, fmt.Errorf("Command definition '%s' already exists", commandDefinition.CommandID)
	}

	return s.data.CommandRepository.Create(commandDefinition)
}

func (s *ManagerState) disableCommandDefinition(commandID string) error {
	return toggleCommandDefinition(s, commandID, false)
}

func (s *ManagerState) enableCommandDefinition(commandID string) error {
	return toggleCommandDefinition(s, commandID, true)
}

func updateCommandDefinition(s *ManagerState, updatedDefinition *CommandDefinition) (*CommandDefinition, error) {
	return s.data.CommandRepository.Update(updatedDefinition)
}

func toggleCommandDefinition(s *ManagerState, commandID string, isEnabled bool) error {
	cd, err := s.getCommandDefinition(commandID)
	if err != nil {
		return err
	}

	cd.Enabled = isEnabled

	_, err = updateCommandDefinition(s, cd)
	if err != nil {
		return err
	}

	err = notifyCommandChange(s)

	return err
}

func notifyCommandChange(s *ManagerState) error {
	commandDefinitions, err := s.getAllCommandDefinitions()
	if err != nil {
		return err
	}

	return s.nats.NotifyCommandConfigChange(commandDefinitions)
}
