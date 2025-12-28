package http

import "github.com/maqsatto/gojson/internal/engine"

type SelectReq struct {
	Table  string             `json:"table"`
	Where  []engine.Condition `json:"where,omitempty"`
	Limit  int                `json:"limit,omitempty"`
	Offset int                `json:"offset,omitempty"`
	SortBy string             `json:"sortBy,omitempty"`
	Desc   bool               `json:"desc,omitempty"`
}

type InsertReq struct {
	Table string     `json:"table"`
	Row   engine.Row `json:"row"`
}

type UpdateReq struct {
	Table string             `json:"table"`
	Where []engine.Condition `json:"where,omitempty"`
	Set   engine.Row         `json:"set"`
}

type DeleteReq struct {
	Table string             `json:"table"`
	Where []engine.Condition `json:"where,omitempty"`
}

type Resp struct {
	OK       bool         `json:"ok"`
	Error    string       `json:"error,omitempty"`
	Result   []engine.Row `json:"result,omitempty"`
	Affected int          `json:"affected,omitempty"`
}
