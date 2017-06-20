package main

import (
	"flag"

	"github.com/he4d/almue/almue"
	_ "github.com/kidoman/embd/host/rpi"
	_ "github.com/mattn/go-sqlite3"
)

var (
	simulate = flag.Bool("simulate", false, "starts simulation mode without gpio (operations will be logged instead)")
	verbose  = flag.Bool("verbose", false, "verbose logging will be enabled")
)

func main() {
	flag.Parse()
	almue := almue.NewAlmue("./almue.db", *simulate, *verbose)
	almue.Serve(":8000")
}
