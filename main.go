package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v2"
)

type Play struct {
	Roles []string `yaml:"roles"`
}

type MetaMain struct {
	Dependencies []struct {
		Role string `yaml:"role"`
	} `yaml:"dependencies"`
}

func getDependencies(roleName string, rolesPath string) ([]string, error) {
	metaFile := filepath.Join(rolesPath, roleName, "meta", "main.yml")
	data, err := os.ReadFile(metaFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var meta MetaMain
	err = yaml.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}

	deps := make([]string, len(meta.Dependencies))
	for i, dep := range meta.Dependencies {
		deps[i] = dep.Role
	}
	return deps, nil
}

func getAllDependencies(roles []string, rolesPath string) (map[string]bool, error) {
	allDeps := make(map[string]bool)
	toProcess := append([]string{}, roles...)

	for len(toProcess) > 0 {
		role := toProcess[0]
		toProcess = toProcess[1:]

		if !allDeps[role] {
			allDeps[role] = true

			deps, err := getDependencies(role, rolesPath)
			if err != nil {
				return nil, fmt.Errorf("error getting dependencies for role %s: %v", role, err)
			}

			toProcess = append(toProcess, deps...)
		}
	}

	return allDeps, nil
}

func getPlaybookRoles(playbookPath string) ([]string, error) {
	file, err := os.Open(playbookPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var plays []Play
	err = yaml.Unmarshal(data, &plays)
	if err != nil {
		return nil, err
	}

	var roles []string
	for _, play := range plays {
		roles = append(roles, play.Roles...)
	}
	return roles, nil
}

func getAllRoles(rolesPath string) ([]string, error) {
	entries, err := os.ReadDir(rolesPath)
	if err != nil {
		return nil, err
	}

	var roles []string
	for _, entry := range entries {
		if entry.IsDir() {
			roles = append(roles, entry.Name())
		}
	}
	return roles, nil
}

func findUnusedRoles(playbookPath, rolesPath string) ([]string, error) {
	playbookRoles, err := getPlaybookRoles(playbookPath)
	if err != nil {
		return nil, fmt.Errorf("error reading playbook: %v", err)
	}

	allUsedRoles, err := getAllDependencies(playbookRoles, rolesPath)
	if err != nil {
		return nil, fmt.Errorf("error getting dependencies: %v", err)
	}

	allRoles, err := getAllRoles(rolesPath)
	if err != nil {
		return nil, fmt.Errorf("error getting all roles: %v", err)
	}

	var unusedRoles []string
	for _, role := range allRoles {
		if !allUsedRoles[role] {
			unusedRoles = append(unusedRoles, role)
		}
	}

	sort.Strings(unusedRoles)
	return unusedRoles, nil
}

func main() {
	playbookPath := flag.String("playbook", "", "Path to the Ansible playbook YAML file")
	rolesPath := flag.String("roles", "", "Path to the Ansible roles directory")
	flag.Parse()

	if *playbookPath == "" || *rolesPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	unusedRoles, err := findUnusedRoles(*playbookPath, *rolesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(unusedRoles) > 0 {
		for _, role := range unusedRoles {
			fmt.Println(role)
		}
	}
}
