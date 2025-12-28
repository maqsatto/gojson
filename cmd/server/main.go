package main

import (
	"github.com/gin-gonic/gin"
	"github.com/maqsatto/gojson/internal/engine"
	"github.com/maqsatto/gojson/internal/http"
)

func main() {
	storage, _ := engine.NewJSONStorage("data/database.json")
	eng := engine.NewJSONEngine(storage)
	r := gin.Default()

	http.Register(r, http.New(eng))

	_ = r.Run(":8080")
}
