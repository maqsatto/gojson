package engine

type Row map[string]any

// Condition Op =, !=, >, <, >=, <=
type Condition struct {
	Field string
	Op    string
	Value any
}

type Query struct {
	Where   []Condition
	OrderBy string
	Limit   int
}
