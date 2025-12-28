package http

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine, h *Handler) {
	api := r.Group("/api")
	{
		api.POST("/select", h.Select)
		api.POST("/insert", h.Insert)
		api.POST("/update", h.Update)
		api.POST("/delete", h.Delete)
	}
}
