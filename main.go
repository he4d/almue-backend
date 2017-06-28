package main

import (
	"flag"

	"github.com/he4d/almue/almue"
	_ "github.com/mattn/go-sqlite3"
)

var (
	simulate  = flag.Bool("simulate", false, "starts simulation mode without gpio (operations will be written to stdout instead)")
	routes    = flag.Bool("routes", false, "generate router documentation")
	publicapi = flag.Bool("publicapi", false, "enables public access to the rest service")
)

func main() {
	flag.Parse()
	almue := almue.NewAlmue("./almue.db", *simulate, *publicapi)
	if *routes {
		almue.GenerateRoutesDoc()
		return
	}
	almue.Serve(":8000")
}
