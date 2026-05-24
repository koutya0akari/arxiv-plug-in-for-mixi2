package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
)

type Store map[string][]string

func Load(path string) (Store, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Store{}, nil
	}
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return Store{}, nil
	}
	var s Store
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	if s == nil {
		s = Store{}
	}
	return s, nil
}

func Save(path string, s Store) error {
	normalized := Store{}
	for category, ids := range s {
		seen := map[string]bool{}
		for _, id := range ids {
			if id != "" {
				seen[id] = true
			}
		}
		for id := range seen {
			normalized[category] = append(normalized[category], id)
		}
		sort.Strings(normalized[category])
	}

	b, err := json.MarshalIndent(normalized, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func (s Store) Has(category, id string) bool {
	for _, existing := range s[category] {
		if existing == id {
			return true
		}
	}
	return false
}

func (s Store) Mark(category, id string) {
	if id == "" || s.Has(category, id) {
		return
	}
	s[category] = append(s[category], id)
}

func (s Store) IsEmpty() bool {
	for _, ids := range s {
		if len(ids) > 0 {
			return false
		}
	}
	return true
}
