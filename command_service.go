package main

import (
	cydata "./lib/data"
)

func (s *managerState) getAllCommandDefinitions() ([]cydata.CommandDefinition, error) {
	commandDefinitions := make([]cydata.CommandDefinition, 0)

	return commandDefinitions, nil
}

func (s *managerState) disableCommandDefinition(commandID string) error {
	cd, err := getCommandDefinition(commandID)
	if err != nil {
		return err
	}

	cd.Enabled = false

	return nil
}

func (s *managerState) enableCommandDefinition(commandID string) error {
	cd, err := getCommandDefinition(commandID)
	if err != nil {
		return err
	}

	cd.Enabled = true

	return nil
}

func getCommandDefinition(commandID string) (cydata.CommandDefinition, error) {
	commandDefinition := cydata.CommandDefinition{}

	return commandDefinition, nil
}
