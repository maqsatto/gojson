package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type JSONStorage struct {
	mu   sync.RWMutex
	path string
	data map[string][]Row
}

func NewJSONStorage(path string) (*JSONStorage, error) {
	s := &JSONStorage{
		path: path,
		data: make(map[string][]Row),
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err := s.saveLockedAtomic(); err != nil {
				return nil, err
			}
			return s, nil
		}
		return nil, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return s, nil
	}
	_ = json.Unmarshal(b, &s.data)
	return s, nil
}

func (s *JSONStorage) saveLockedAtomic() error {
	tmp := s.path + ".tmp"

	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
