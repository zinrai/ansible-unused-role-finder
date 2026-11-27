package main

import (
	"path/filepath"
	"testing"
)

func TestExpandRolesTransitiveWithCycle(t *testing.T) {
	roles := t.TempDir()
	// a depends on b, b depends on c, c depends back on a (cycle).
	writeFile(t, filepath.Join(roles, "a", "meta", "main.yml"), "dependencies:\n  - role: b\n")
	writeFile(t, filepath.Join(roles, "b", "meta", "main.yml"), "dependencies:\n  - role: c\n")
	writeFile(t, filepath.Join(roles, "c", "meta", "main.yml"), "dependencies:\n  - role: a\n")
	// d has no meta file at all.

	used, err := expandRoles([]string{"a", "d"}, roles)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range []string{"a", "b", "c", "d"} {
		if !used[r] {
			t.Fatalf("expected %q to be reached, got %v", r, used)
		}
	}
}

func TestExpandRolesRejectsUnresolvableDependency(t *testing.T) {
	entries := map[string]string{
		"name key":    "- name: b",
		"no role key": "- tags: web",
		"path":        "- role: roles/b",
		"template":    `- role: "{{ b }}"`,
	}

	for name, entry := range entries {
		t.Run(name, func(t *testing.T) {
			roles := t.TempDir()
			writeFile(t, filepath.Join(roles, "a", "meta", "main.yml"), "dependencies:\n  "+entry+"\n")

			used, err := expandRoles([]string{"a"}, roles)
			if err == nil {
				t.Fatalf("got %v, want an error", used)
			}
		})
	}
}
