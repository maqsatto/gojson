package gojson

import (
	"github.com/gin-gonic/gin"

	"github.com/maqsatto/gojson/internal/engine"
	"github.com/maqsatto/gojson/internal/http"
)

type Server struct {
	Router *gin.Engine
}

func New(path string) (*Server, error) {
	st, err := engine.NewJSONStorage(path)
	if err != nil {
		return nil, err
	}
	eng := engine.NewJSONEngine(st)

	r := gin.Default()
	http.Register(r, http.New(eng))

	return &Server{Router: r}, nil
}

func (s *Server) Run(addr string) error {
	return s.Router.Run(addr)
}
