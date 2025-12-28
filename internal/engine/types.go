package engine

type Row map[string]any

// Condition Op =, !=, >, <, >=, <=
type Condition struct {
	Field string `json:"field"`
	Op    string `json:"op"`
	Value any    `json:"value"`
}

type Query struct {
	Where  []Condition `json:"where,omitempty"`
	Limit  int         `json:"limit,omitempty"`
	Offset int         `json:"offset,omitempty"`
	SortBy string      `json:"sortBy,omitempty"`
	Desc   bool        `json:"desc,omitempty"`
}
