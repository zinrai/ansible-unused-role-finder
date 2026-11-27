package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type roleRef struct {
	Name string
}

func (r *roleRef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var entry interface{}
	if err := unmarshal(&entry); err != nil {
		return err
	}

	name, ok := entry.(string)
	if mapping, isMapping := entry.(map[string]interface{}); isMapping {
		name, ok = mapping["role"].(string)
	}
	if !ok {
		return fmt.Errorf("role entry %v is neither a role name nor a mapping with a role key", entry)
	}
	r.Name = name
	return nil
}

type play struct {
	Roles []roleRef `yaml:"roles"`
}

func playbookRoles(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var plays []play
	if err := yaml.Unmarshal(data, &plays); err != nil {
		return nil, err
	}

	var roles []string
	for _, p := range plays {
		names, err := roleNames(p.Roles)
		if err != nil {
			return nil, err
		}
		roles = append(roles, names...)
	}
	return roles, nil
}

func roleNames(refs []roleRef) ([]string, error) {
	var names []string
	for _, r := range refs {
		if !isRoleDirName(r.Name) {
			return nil, fmt.Errorf("role %q is not the name of a directory in the roles directory", r.Name)
		}
		names = append(names, r.Name)
	}
	return names, nil
}

func isRoleDirName(name string) bool {
	return name != "" && !strings.Contains(name, "/") && !strings.Contains(name, "{{")
}
