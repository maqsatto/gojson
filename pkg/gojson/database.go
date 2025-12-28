package gojson

import (
	"github.com/maqsatto/gojson/internal/engine"
	"github.com/maqsatto/gojson/internal/query"
)

type DB struct {
	eng  engine.Engine
	exec *query.Executor
}

func Open(path string) (*DB, error) {
	st, err := engine.NewJSONStorage(path)
	if err != nil {
		return nil, err
	}
	eng := engine.NewJSONEngine(st)

	return &DB{
		eng:  eng,
		exec: query.NewExecutor(eng),
	}, nil
}

func (db *DB) From(table string) *query.Builder {
	return query.From(table)
}

func (db *DB) All(b *query.Builder) ([]engine.Row, error) {
	return db.exec.All(b)
}

func (db *DB) Insert(table string, row map[string]any) error {
	return db.eng.Insert(table, row)
}

func (db *DB) Update(b *query.Builder, set map[string]any) (int, error) {
	return db.exec.Update(b, set)
}

func (db *DB) Delete(b *query.Builder) (int, error) {
	return db.exec.Delete(b)
}
