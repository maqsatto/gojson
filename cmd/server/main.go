package main

import (
	"log"

	"github.com/maqsatto/gojson/pkg/gojson"
)

func main() {
	srv, err := gojson.New("data/db.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(srv.Run(":8080"))
}
