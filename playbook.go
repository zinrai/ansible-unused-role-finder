package main

import (
	"os"

	"github.com/goccy/go-yaml"
)

// roleRef accepts both forms a play may use in its roles: list,
// a bare "rolename" and a "{ role: rolename }" mapping.
type roleRef struct {
	Name string
}

func (r *roleRef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var name string
	if err := unmarshal(&name); err == nil {
		r.Name = name
		return nil
	}

	var mapping struct {
		Role string `yaml:"role"`
	}
	if err := unmarshal(&mapping); err != nil {
		return err
	}
	r.Name = mapping.Role
	return nil
}

type play struct {
	Roles []roleRef `yaml:"roles"`
}

// playbookRoles returns the roles declared at play level across every play in a
// playbook file.
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
		roles = append(roles, roleNames(p.Roles)...)
	}
	return roles, nil
}

// roleNames extracts the non-empty role names from a play's roles list.
func roleNames(refs []roleRef) []string {
	var names []string
	for _, r := range refs {
		if r.Name != "" {
			names = append(names, r.Name)
		}
	}
	return names
}
