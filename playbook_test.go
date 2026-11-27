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
    - nginx
    - role: docker.io
- hosts: db
  roles:
    - postgresql
`)

	got, err := playbookRoles(path)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"nginx", "docker.io", "postgresql"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestPlaybookRolesRejectsOtherForms(t *testing.T) {
	entries := map[string]string{
		"name key":      "- name: nginx",
		"no role key":   "- tags: web",
		"path":          "- roles/nginx",
		"template":      `- "{{ webserver }}"`,
		"number":        "- 123",
		"null":          "- ~",
		"empty name":    `- role: ""`,
		"non-string":    "- role: 123",
		"nested string": "- [nginx]",
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
