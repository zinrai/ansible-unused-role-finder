package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// A role unused by one playbook but used by another that shares the roles
// directory must not be reported as unused.
func TestUnusedRolesSharedAcrossPlaybooks(t *testing.T) {
	roles := t.TempDir()
	for _, r := range []string{"role1", "role2", "role3", "orphan"} {
		writeFile(t, filepath.Join(roles, r, "tasks", "main.yml"), "- name: noop\n  debug:\n    msg: ok\n")
	}
	// role2 depends on role3 through meta.
	writeFile(t, filepath.Join(roles, "role2", "meta", "main.yml"), "dependencies:\n  - role: role3\n")

	dir := t.TempDir()
	playbook1 := filepath.Join(dir, "playbook1.yml")
	playbook2 := filepath.Join(dir, "playbook2.yml")
	writeFile(t, playbook1, "- hosts: all\n  roles:\n    - role1\n")
	writeFile(t, playbook2, "- hosts: all\n  roles:\n    - role2\n")

	got, err := unusedRoles([]string{playbook1, playbook2}, roles)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"orphan"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestUnusedRolesListsSymlinkedRole(t *testing.T) {
	roles := t.TempDir()
	target := t.TempDir()
	writeFile(t, filepath.Join(target, "tasks", "main.yml"), "- name: noop\n  debug:\n    msg: ok\n")
	if err := os.Symlink(target, filepath.Join(roles, "linked")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(target, "absent"), filepath.Join(roles, "dangling")); err != nil {
		t.Fatal(err)
	}

	playbook := filepath.Join(t.TempDir(), "playbook.yml")
	writeFile(t, playbook, "- hosts: all\n  roles: []\n")

	got, err := unusedRoles([]string{playbook}, roles)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"linked"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
