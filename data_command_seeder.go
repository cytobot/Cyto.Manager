package main

import (
	"time"
)

func (r *CommandRepository) seedInitialCommandDefinitions() error {
	storedDefinitions, err := r.GetAll()
	if err != nil {
		return err
	}

	for _, initialDefinition := range getInitialCommandDefinitions() {
		if containsCommandDefinition(initialDefinition, storedDefinitions) {
			continue
		}

		_, err := r.Create(initialDefinition)
		if err != nil {
			return err
		}
	}

	return nil
}

func containsCommandDefinition(targetDefinition *CommandDefinition, sourceDefinitions []*CommandDefinition) bool {
	if targetDefinition == nil || len(sourceDefinitions) == 0 {
		return false
	}

	for _, sourceDefinition := range sourceDefinitions {
		if targetDefinition.CommandID == sourceDefinition.CommandID {
			return true
		}
	}

	return false
}

func getInitialCommandDefinitions() []*CommandDefinition {
	return []*CommandDefinition{
		&CommandDefinition{
			CommandID:            "invite",
			Enabled:              true,
			Unlisted:             false,
			Description:          "Get an invite link to add this bot to your server!",
			Triggers:             []string{"invite"},
			PermissionLevel:      "USER",
			ParameterDefinitions: nil,
			LastModifiedDateUtc:  time.Now().UTC(),
			LastModifiedUserID:   "",
		},
		&CommandDefinition{
			CommandID:            "stats",
			Enabled:              true,
			Unlisted:             false,
			Description:          "Get bot stats",
			Triggers:             []string{"stats"},
			PermissionLevel:      "USER",
			ParameterDefinitions: nil,
			LastModifiedDateUtc:  time.Now().UTC(),
			LastModifiedUserID:   "",
		},
		&CommandDefinition{
			CommandID:            "stats-plus",
			Enabled:              true,
			Unlisted:             true,
			Description:          "Get bot stats",
			Triggers:             []string{"stats+"},
			PermissionLevel:      "OWNER",
			ParameterDefinitions: nil,
			LastModifiedDateUtc:  time.Now().UTC(),
			LastModifiedUserID:   "",
		},
	}
}
