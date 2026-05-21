package types

import (
	"fmt"
	"os"
	"sigs.k8s.io/yaml"
)

func ReadFileAsString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ParseConfigDefinition(path string) (*LiveKitConfigDefinition, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Attempt to unmarshall the contents
	var def *LiveKitConfigDefinition
	err = yaml.Unmarshal(content, &def)
	if err != nil {
		return nil, err
	}
	// Populate secrets from files
	if e := loadPasswordFiles(def); e != nil {
		return nil, fmt.Errorf("Failed to load password from file: %s", e)
	}
	// Perform validation
	if e := validateTrunks(def); e != nil {
		return nil, fmt.Errorf("Failed to validate trunks: %s", e)
	}
	if e := validateDispatchRules(def); e != nil {
		return nil, fmt.Errorf("Failed to validate Dispatch Rules: %s", e)
	}
	return def, nil
}

func empty(s string) bool {
	return len(s) == 0
}

func loadPasswordFiles(def *LiveKitConfigDefinition) error {
	if def.SIPConfig == nil || def.SIPConfig.Trunks == nil {
		return nil
	}
	inbound := def.SIPConfig.Trunks.InboundTrunks
	for idx := range inbound {
		trunk := inbound[idx]
		if !empty(trunk.PasswordFile) {
			password, err := os.ReadFile(trunk.PasswordFile)
			if err != nil {
				return fmt.Errorf("Failed to load password file (%s): %s", trunk.PasswordFile, err)
			}
			trunk.Password = string(password)
		}
	}
	outbound := def.SIPConfig.Trunks.OutboundTrunks
	for idx := range outbound {
		trunk := outbound[idx]
		if !empty(trunk.PasswordFile) {
			password, err := os.ReadFile(trunk.PasswordFile)
			if err != nil {
				return fmt.Errorf("Failed to load password file (%s): %s", trunk.PasswordFile, err)
			}
			trunk.Password = string(password)
		}
	}
	return nil
}

func validateTrunks(def *LiveKitConfigDefinition) error {
	if def.SIPConfig == nil || def.SIPConfig.Trunks == nil {
		return nil
	}
	inbound := def.SIPConfig.Trunks.InboundTrunks
	for idx := range inbound {
		trunk := inbound[idx]
		if empty(trunk.Name) {
			return fmt.Errorf("all inbound trunks must have a non empty name")
		}
		if !empty(trunk.Username) && empty(trunk.Password) {
			return fmt.Errorf("username provided for trunk (%s) but no password provided", trunk.Name)
		}
	}
	outbound := def.SIPConfig.Trunks.OutboundTrunks
	for idx := range outbound {
		trunk := outbound[idx]
		if empty(trunk.Name) {
			return fmt.Errorf("all outbound trunks must have a non empty name")
		}
		if empty(trunk.Address) {
			return fmt.Errorf("all outbound trunks must have a non empty address")
		}
		if !empty(trunk.Username) && empty(trunk.Password) {
			return fmt.Errorf("username provided for trunk (%s) but no password provided", trunk.Name)
		}
	}
	return nil
}

func validateDispatchRules(def *LiveKitConfigDefinition) error {
	if def.SIPConfig == nil {
		return nil
	}
	rules := def.SIPConfig.DispatchRules
	for idx := range rules {
		rule := rules[idx]
		if empty(rule.Name) {
			return fmt.Errorf("all dispatch rules must have a non empty name")
		}
		numRules := 0
		if rule.Rule.DispatchRuleDirect != nil {
			numRules += 1
		}
		if rule.Rule.DispatchRuleIndividual != nil {
			numRules += 1
		}
		if rule.Rule.DispatchRuleCallee != nil {
			numRules += 1
		}
		if numRules != 1 {
			return fmt.Errorf(
				"a dispatch rule (%s) must have one, and only one rule defined (has %d)",
				rule.Name,
				numRules)
		}
	}
	return nil
}
