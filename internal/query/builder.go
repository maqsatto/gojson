package query

import "github.com/maqsatto/gojson/internal/engine"

type Builder struct {
	table string
	query engine.Query
}

func From(table string) *Builder {
	return &Builder{table: table}
}

func (b *Builder) Where(field, op string, value any) *Builder {
	b.query.Where = append(b.query.Where, engine.Condition{
		Field: field,
		Op:    op,
		Value: value,
	})
	return b
}
