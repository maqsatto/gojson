package engine

type Row map[string]any

// Condition Op =, !=, >, <, >=, <=
type Condition struct {
	Field string `json:"field"`
	Op    string `json:"op"`
	Value any    `json:"value"`
}

type Expr struct {
	// ровно один из этих вариантов используется:
	Cond  *Condition `json:"cond,omitempty"`
	Group *Group     `json:"group,omitempty"`
}

type Group struct {
	Op   string `json:"op"` // "AND" | "OR"
	Expr []Expr `json:"expr"`
}
type Query struct {
	Where  []Condition `json:"where,omitempty"`
	Filter *Expr       `json:"filter,omitempty"`
	Limit  int         `json:"limit,omitempty"`
	Offset int         `json:"offset,omitempty"`
	SortBy string      `json:"sortBy,omitempty"`
	Desc   bool        `json:"desc,omitempty"`
}
