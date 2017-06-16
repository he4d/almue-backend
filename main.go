package main

import (
	"flag"

	"github.com/he4d/almue/almue"
	_ "github.com/kidoman/embd/host/rpi"
	_ "github.com/mattn/go-sqlite3"
)

var (
	simulate = flag.Bool("simulate", false, "starts simulation mode without gpio (operations will be logged instead)")
)

func main() {
	flag.Parse()
	almue := almue.Almue{}
	almue.Initialize("./almue.db", *simulate)
	almue.Run(":8000")
}
