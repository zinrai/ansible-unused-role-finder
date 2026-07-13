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
