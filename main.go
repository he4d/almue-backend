package main

import (
	"flag"

	"github.com/he4d/almue/almue"
	// _ "github.com/kidoman/embd/host/rpi"
	_ "github.com/mattn/go-sqlite3"
)

var (
	simulate = flag.Bool("simulate", false, "starts simulation mode without gpio (operations will be logged instead)")
	routes   = flag.Bool("routes", false, "Generate router documentation")
)

func main() {
	flag.Parse()
	almue := almue.NewAlmue("./almue.db", *simulate)
	if *routes {
		almue.GenerateRoutesDoc()
		return
	}
	almue.Serve(":8000")
}
