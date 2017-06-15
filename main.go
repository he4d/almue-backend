package main

import (
	"github.com/he4d/almue/api"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	api := api.API{}
	api.Initialize("./almue.db")
	api.Run(":8000")
}
