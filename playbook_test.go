package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestPlaybookRolesBothForms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "playbook.yml")
	writeFile(t, path, `
- hosts: all
  roles:
    - role1
    - role: dotted.role
- hosts: db
  roles:
    - role2
`)

	got, err := playbookRoles(path)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"role1", "dotted.role", "role2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestPlaybookRolesRejectsOtherForms(t *testing.T) {
	entries := map[string]string{
		"name key":      "- name: role1",
		"no role key":   "- tags: tag1",
		"path":          "- roles/role1",
		"template":      `- "{{ role1 }}"`,
		"number":        "- 123",
		"null":          "- ~",
		"empty name":    `- role: ""`,
		"non-string":    "- role: 123",
		"nested string": "- [role1]",
	}

	for name, entry := range entries {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "playbook.yml")
			writeFile(t, path, "- hosts: all\n  roles:\n    "+entry+"\n")

			got, err := playbookRoles(path)
			if err == nil {
				t.Fatalf("got %v, want an error", got)
			}
		})
	}
}
