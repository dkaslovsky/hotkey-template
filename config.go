package main

import (
	"encoding/json"
	"fmt"
	"io"

	"golang.design/x/hotkey"
)

// Configuration for each hotkey trigger
type keyConfig struct {
	Name         string   `json:"name"`
	CommandName  string   `json:"command_name"`
	CommandArgs  []string `json:"command_args"`
	Key          string   `json:"key"`
	KeyModifiers []string `json:"key_modifiers"`
}

func parseKeyConfig(r io.Reader) ([]*hotkeyTrigger, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	keyConfigs := []*keyConfig{}
	err = json.Unmarshal(raw, &keyConfigs)
	if err != nil {
		return nil, err
	}

	hotkeyTriggers := []*hotkeyTrigger{}
	for _, config := range keyConfigs {

		// Convert string representation of keys and modifiers to internal constants
		key, ok := keyMap[config.Key]
		if !ok {
			return nil, fmt.Errorf("unknown key: %s", config.Key)
		}
		modifiers := []hotkey.Modifier{}
		for _, mod := range config.KeyModifiers {
			modifier, ok := modifierMap[mod]
			if !ok {
				return nil, fmt.Errorf("unknown key modifier: %s", mod)
			}
			modifiers = append(modifiers, modifier)
		}

		hotkeyTriggers = append(hotkeyTriggers, &hotkeyTrigger{
			Hotkey:      hotkey.New(modifiers, key),
			name:        config.Name,
			commandName: config.CommandName,
			commandArgs: config.CommandArgs,
		})
	}

	return hotkeyTriggers, nil
}
