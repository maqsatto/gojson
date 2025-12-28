package query

import "github.com/maqsatto/gojson/internal/engine"

type Executor struct {
	engine engine.Engine
}

func NewExecutor(e engine.Engine) *Executor {
	return &Executor{engine: e}
}

func (e Executor) All(b *Builder) ([]engine.Row, error) {
	return e.engine.Select(b.table, b.query)
}

func (e Executor) Insert(table string, row engine.Row) error {
	return e.engine.Insert(table, row)
}
