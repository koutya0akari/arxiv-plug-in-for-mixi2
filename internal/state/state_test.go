package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingReturnsEmptyStore(t *testing.T) {
	s, err := Load(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !s.IsEmpty() {
		t.Fatalf("Load() = %+v, want empty", s)
	}
}

func TestSaveLoadAndMark(t *testing.T) {
	path := filepath.Join(t.TempDir(), "posted.json")
	s := Store{}
	s.Mark("math.CT", "2501.00002")
	s.Mark("math.CT", "2501.00001")
	s.Mark("math.CT", "2501.00001")

	if err := Save(path, s); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !loaded.Has("math.CT", "2501.00001") || !loaded.Has("math.CT", "2501.00002") {
		t.Fatalf("loaded state missing ids: %+v", loaded)
	}
	if got, want := len(loaded["math.CT"]), 2; got != want {
		t.Fatalf("len(ids) = %d, want %d", got, want)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(b) == "" || b[len(b)-1] != '\n' {
		t.Fatalf("saved JSON should end with newline: %q", b)
	}
}
