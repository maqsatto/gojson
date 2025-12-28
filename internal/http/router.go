package http

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine, h *Handler) {
	api := r.Group("/api")
	api.POST("/insert", h.Insert)
}
