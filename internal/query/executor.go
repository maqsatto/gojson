package query

import "github.com/maqsatto/gojson/internal/engine"

type Executor struct {
	eng engine.Engine
}

func NewExecutor(e engine.Engine) *Executor {
	return &Executor{eng: e}
}

func (x *Executor) All(b *Builder) ([]engine.Row, error) {
	return x.eng.Select(b.Table(), b.Query())
}

func (x *Executor) Insert(table string, row engine.Row) error {
	return x.eng.Insert(table, row)
}

func (x *Executor) Update(b *Builder, set engine.Row) (int, error) {
	return x.eng.Update(b.Table(), b.Query(), set)
}

func (x *Executor) Delete(b *Builder) (int, error) {
	return x.eng.Delete(b.Table(), b.Query())
}
