package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

// A role unused by one playbook but used by another that shares the roles
// directory must not be reported as unused.
func TestUnusedRolesSharedAcrossPlaybooks(t *testing.T) {
	roles := t.TempDir()
	for _, r := range []string{"consul", "nomad", "docker.io", "orphan"} {
		writeFile(t, filepath.Join(roles, r, "tasks", "main.yml"), "- name: noop\n  debug:\n    msg: ok\n")
	}
	// nomad depends on docker.io through meta.
	writeFile(t, filepath.Join(roles, "nomad", "meta", "main.yml"), "dependencies:\n  - role: docker.io\n")

	dir := t.TempDir()
	consulPB := filepath.Join(dir, "consul.yml")
	nomadPB := filepath.Join(dir, "nomad.yml")
	writeFile(t, consulPB, "- hosts: all\n  roles:\n    - consul\n")
	writeFile(t, nomadPB, "- hosts: all\n  roles:\n    - nomad\n")

	// Checked alone, consul.yml would wrongly flag nomad and docker.io.
	got, err := unusedRoles([]string{consulPB, nomadPB}, roles)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"orphan"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
