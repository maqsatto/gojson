package query

import "github.com/maqsatto/gojson/internal/engine"

type Builder struct {
	table string
	q     engine.Query
}

func From(table string) *Builder {
	return &Builder{table: table}
}

func (b *Builder) Where(field, op string, value any) *Builder {
	b.q.Where = append(b.q.Where, engine.Condition{
		Field: field,
		Op:    op,
		Value: value,
	})
	return b
}

func (b *Builder) Limit(n int) *Builder {
	b.q.Limit = n
	return b
}

func (b *Builder) Offset(n int) *Builder {
	b.q.Offset = n
	return b
}

func (b *Builder) SortBy(field string, desc bool) *Builder {
	b.q.SortBy = field
	b.q.Desc = desc
	return b
}

func (b *Builder) Table() string       { return b.table }
func (b *Builder) Query() engine.Query { return b.q }
