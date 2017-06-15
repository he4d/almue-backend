package main

import (
	"os"

	"github.com/he4d/almue/almue"
	_ "github.com/kidoman/embd/host/rpi"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	apiOnly := false
	if len(os.Args) == 2 && os.Args[1] == "--api-only" {
		apiOnly = true
	}
	almue := almue.Almue{}
	almue.Initialize("./almue.db", apiOnly)
	almue.Run(":8000")
}
