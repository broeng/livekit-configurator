package types

import (
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

	return def, nil
}
