package gojson

import (
	"github.com/maqsatto/gojson/internal/engine"
	"github.com/maqsatto/gojson/internal/query"
)

type Client struct {
	exec *query.Executor
}

func New(eng engine.Engine) *Client {
	return &Client{exec: query.NewExecutor(eng)}
}

func (c *Client) From(table string) *query.Builder {
	return query.From(table)
}

func (c *Client) All(b *query.Builder) ([]engine.Row, error) {
	return c.exec.All(b)
}

func (c *Client) Insert(table string, row engine.Row) error {
	return c.exec.Insert(table, row)
}

func (c *Client) Update(b *query.Builder, set engine.Row) (int, error) {
	return c.exec.Update(b, set)
}

func (c *Client) Delete(b *query.Builder) (int, error) {
	return c.exec.Delete(b)
}
