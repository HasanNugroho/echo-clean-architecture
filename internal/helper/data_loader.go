package helper

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadStringListFromYAML(path string, key string) (map[string]struct{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var parsed map[string][]string
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	permSet := make(map[string]struct{})
	for _, perm := range parsed[key] {
		permSet[perm] = struct{}{}
	}

	return permSet, nil
}
