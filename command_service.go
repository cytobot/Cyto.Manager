package main

func (s *managerState) getAllCommandDefinitions() ([]*CommandDefinition, error) {
	commandDefinitions, err := s.data.CommandRepository.GetAll()

	return commandDefinitions, err
}

func (s *managerState) getCommandDefinition(commandID string) (*CommandDefinition, error) {
	commandDefinition, err := s.data.CommandRepository.Get(commandID)

	return commandDefinition, err
}

func (s *managerState) disableCommandDefinition(commandID string) error {
	return toggleCommandDefinition(s, commandID, false)
}

func (s *managerState) enableCommandDefinition(commandID string) error {
	return toggleCommandDefinition(s, commandID, true)
}

func updateCommandDefinition(s *managerState, updatedDefinition *CommandDefinition) (*CommandDefinition, error) {
	return s.data.CommandRepository.Update(updatedDefinition)
}

func toggleCommandDefinition(s *managerState, commandID string, isEnabled bool) error {
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

func notifyCommandChange(s *managerState) error {
	commandDefinitions, err := s.getAllCommandDefinitions()
	if err != nil {
		return err
	}

	return s.nats.notifyCommandConfigChange(commandDefinitions)
}
