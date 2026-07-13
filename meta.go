package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type roleMeta struct {
	Dependencies []struct {
		Role string `yaml:"role"`
	} `yaml:"dependencies"`
}

// roleDependencies returns the roles listed in roles/<role>/meta/main.yml.
// A role without a meta file has no dependencies.
func roleDependencies(role, rolesPath string) ([]string, error) {
	metaFile := filepath.Join(rolesPath, role, "meta", "main.yml")
	data, err := os.ReadFile(metaFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var meta roleMeta
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	var deps []string
	for _, d := range meta.Dependencies {
		if d.Role != "" {
			deps = append(deps, d.Role)
		}
	}
	return deps, nil
}

// expandRoles walks the meta-dependency graph from the seed roles and returns
// the full set of roles reached. Cycles terminate on the visited set.
func expandRoles(seed []string, rolesPath string) (map[string]bool, error) {
	used := make(map[string]bool)
	queue := append([]string{}, seed...)

	for len(queue) > 0 {
		role := queue[0]
		queue = queue[1:]
		if used[role] {
			continue
		}
		used[role] = true

		deps, err := roleDependencies(role, rolesPath)
		if err != nil {
			return nil, fmt.Errorf("role %s: %w", role, err)
		}
		queue = append(queue, deps...)
	}
	return used, nil
}
