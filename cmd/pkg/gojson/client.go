package gojson

import (
	"github.com/maqsatto/gojson/internal/engine"
	"github.com/maqsatto/gojson/internal/query"
)

type Client struct {
	exec *query.Executor
}

func New(e engine.Engine) *Client {
	return &Client{exec: query.NewExecutor(e)}
}

func (c *Client) From(table string) *query.Builder {
	return query.From(table)
}

func (c *Client) Insert(table string, row engine.Row) error {
	return c.exec.Insert(table, row)
}

func (c *Client) All(b *query.Builder) (rows []engine.Row, err error) {
	return c.exec.All(b)
}
