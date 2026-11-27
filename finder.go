package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func allRoles(rolesPath string) ([]string, error) {
	entries, err := os.ReadDir(rolesPath)
	if err != nil {
		return nil, err
	}

	var roles []string
	for _, e := range entries {
		if isDir(filepath.Join(rolesPath, e.Name())) {
			roles = append(roles, e.Name())
		}
	}
	return roles, nil
}

// Not DirEntry.IsDir: it reports false for a symlink to a directory, which
// would leave a symlinked role out of the listing.
func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

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

func reachedFromPlaybook(playbook, rolesPath string) (map[string]bool, error) {
	seed, err := playbookRoles(playbook)
	if err != nil {
		return nil, fmt.Errorf("playbook %s: %w", playbook, err)
	}
	return expandRoles(seed, rolesPath)
}

func addAll(dst, src map[string]bool) {
	for k := range src {
		dst[k] = true
	}
}
