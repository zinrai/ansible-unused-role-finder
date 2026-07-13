package main

import (
	"fmt"
	"os"
	"sort"
)

// allRoles lists the role directories under rolesPath.
func allRoles(rolesPath string) ([]string, error) {
	entries, err := os.ReadDir(rolesPath)
	if err != nil {
		return nil, err
	}

	var roles []string
	for _, e := range entries {
		if e.IsDir() {
			roles = append(roles, e.Name())
		}
	}
	return roles, nil
}

// unusedRoles returns the roles under rolesPath that no playbook reaches, either
// directly or through meta dependencies. A role counts as used when any of the
// playbooks reaches it, so playbooks sharing a roles directory must be passed
// together.
func unusedRoles(playbooks []string, rolesPath string) ([]string, error) {
	used := make(map[string]bool)
	for _, pb := range playbooks {
		reached, err := reachedFromPlaybook(pb, rolesPath)
		if err != nil {
			return nil, err
		}
		addAll(used, reached)
	}

	roles, err := allRoles(rolesPath)
	if err != nil {
		return nil, err
	}

	var unused []string
	for _, r := range roles {
		if !used[r] {
			unused = append(unused, r)
		}
	}
	sort.Strings(unused)
	return unused, nil
}

// reachedFromPlaybook returns the roles a single playbook reaches, directly or
// through meta dependencies.
func reachedFromPlaybook(playbook, rolesPath string) (map[string]bool, error) {
	seed, err := playbookRoles(playbook)
	if err != nil {
		return nil, fmt.Errorf("playbook %s: %w", playbook, err)
	}
	return expandRoles(seed, rolesPath)
}

// addAll marks every key of src as present in dst.
func addAll(dst, src map[string]bool) {
	for k := range src {
		dst[k] = true
	}
}
