package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/maqsatto/gojson/internal/engine"
	"github.com/maqsatto/gojson/internal/http"
)

func main() {
	st, err := engine.NewJSONStorage("data/database.json")
	if err != nil {
		log.Fatal(err)
	}
	eng := engine.NewJSONEngine(st)

	r := gin.Default()
	http.Register(r, http.New(eng))

	log.Println("server listening on :8080")
	_ = r.Run(":8080")
}
