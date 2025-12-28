package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maqsatto/gojson/internal/engine"
)

type Handler struct {
	eng engine.Engine
}

func New(eng engine.Engine) *Handler {
	return &Handler{eng: eng}
}

func (h *Handler) Select(c *gin.Context) {
	var req SelectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Resp{OK: false, Error: err.Error()})
		return
	}

	rows, err := h.eng.Select(req.Table, engine.Query{
		Where:  req.Where,
		Limit:  req.Limit,
		Offset: req.Offset,
		SortBy: req.SortBy,
		Desc:   req.Desc,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, Resp{OK: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Resp{OK: true, Result: rows})
}

func (h *Handler) Insert(c *gin.Context) {
	var req InsertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Resp{OK: false, Error: err.Error()})
		return
	}

	if err := h.eng.Insert(req.Table, req.Row); err != nil {
		c.JSON(http.StatusBadRequest, Resp{OK: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Resp{OK: true, Affected: 1})
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Resp{OK: false, Error: err.Error()})
		return
	}

	affected, err := h.eng.Update(req.Table, engine.Query{Where: req.Where}, req.Set)
	if err != nil {
		c.JSON(http.StatusBadRequest, Resp{OK: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Resp{OK: true, Affected: affected})
}

func (h *Handler) Delete(c *gin.Context) {
	var req DeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Resp{OK: false, Error: err.Error()})
		return
	}

	affected, err := h.eng.Delete(req.Table, engine.Query{Where: req.Where})
	if err != nil {
		c.JSON(http.StatusBadRequest, Resp{OK: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Resp{OK: true, Affected: affected})
}
