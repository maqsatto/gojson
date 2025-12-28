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

	idIndex map[string]map[float64]int
	maxID   map[string]float64
}

func NewJSONStorage(path string) (*JSONStorage, error) {
	s := &JSONStorage{
		path:    path,
		data:    make(map[string][]Row),
		idIndex: make(map[string]map[float64]int),
		maxID:   make(map[string]float64),
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
	if len(b) > 0 {
		_ = json.Unmarshal(b, &s.data)
	}

	s.rebuildIndexLocked()
	return s, nil
}

func (s *JSONStorage) rebuildIndexLocked() {
	s.idIndex = make(map[string]map[float64]int)
	s.maxID = make(map[string]float64)

	for table, rows := range s.data {
		s.idIndex[table] = make(map[float64]int)
		var mx float64 = 0
		for i, r := range rows {
			if idv, ok := r["id"]; ok {
				if idf, ok := toFloat(idv); ok {
					s.idIndex[table][idf] = i
					if idf > mx {
						mx = idf
					}
				}
			}
		}
		s.maxID[table] = mx
	}
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
