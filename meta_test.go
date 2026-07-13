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
