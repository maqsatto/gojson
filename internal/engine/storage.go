package engine

import (
	"encoding/json"
	"os"
	"sync"
)

type JSONStorage struct {
	mu   sync.RWMutex
	path string
	data map[string][]Row
}

func NewJSONStorage(path string) (*JSONStorage, error) {
	storage := &JSONStorage{
		path: path,
		data: make(map[string][]Row),
	}
	if b, err := os.ReadFile(path); err == nil && len(b) > 0 {
		_ = json.Unmarshal(b, &storage.data)
	}
	return storage, nil
}

func (s *JSONStorage) save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, _ := json.MarshalIndent(s.data, "", "  ")
	return os.WriteFile(s.path, b, 0644)
}
