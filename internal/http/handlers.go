package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maqsatto/gojson/internal/engine"
)

type Handler struct {
	engine engine.Engine
}

func New(engine engine.Engine) *Handler {
	return &Handler{engine: engine}
}

func (h *Handler) Insert(c *gin.Context) {
	var request struct {
		Table string     `json:"table"`
		Row   engine.Row `json:"row"`
	}
	_ = c.BindJSON(&request)
	_ = h.engine.Insert(request.Table, request.Row)
	c.Status(http.StatusOK)

}
