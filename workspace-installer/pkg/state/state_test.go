package state

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vpsik/workspace-installer/pkg/detector"
)

func TestBuild(t *testing.T) {
	detection := &detector.Result{
		Services: []detector.ServiceInfo{
			{Name: "gitea", Status: detector.StatusInstalled, Details: "Running"},
			{Name: "ollama", Status: detector.StatusMissing, Details: "Not found"},
		},
	}

	s := Build(detection)
	if len(s.Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(s.Services))
	}

	if s.Services[0].Name != "gitea" || s.Services[0].Status != detector.StatusInstalled {
		t.Error("gitea should be installed")
	}
	if s.Services[1].Name != "ollama" || s.Services[1].Status != detector.StatusMissing {
		t.Error("ollama should be missing")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	original := &State{
		Services: []ServiceState{
			{Name: "gitea", Status: detector.StatusInstalled, Details: "Running"},
		},
	}

	if err := original.Save(path); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := LoadState(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if len(loaded.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(loaded.Services))
	}
	if loaded.Services[0].Name != "gitea" {
		t.Errorf("expected gitea, got %s", loaded.Services[0].Name)
	}
}

func TestLoadState_FileNotFound(t *testing.T) {
	_, err := LoadState("/tmp/nonexistent-state.json")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestLoadState_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("{invalid}"), 0644)

	_, err := LoadState(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestBuild_Empty(t *testing.T) {
	s := Build(&detector.Result{})
	if len(s.Services) != 0 {
		t.Errorf("expected 0 services, got %d", len(s.Services))
	}
}
