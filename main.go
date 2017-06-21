package main

import (
	"flag"
	"log"
	"os"

	"github.com/he4d/almue/almue"
	_ "github.com/kidoman/embd/host/rpi"
	_ "github.com/mattn/go-sqlite3"
)

var (
	simulate = flag.Bool("simulate", false, "starts simulation mode without gpio (operations will be logged instead)")
	verbose  = flag.Bool("verbose", false, "verbose logging will be enabled")
	routes   = flag.Bool("routes", false, "Generate router documentation")
)

func main() {
	flag.Parse()
	almue := almue.NewAlmue("./almue.db", *simulate, *verbose)
	// Create only the Routes document, then exit
	if *routes {
		content := almue.GenerateRoutes()
		f, err := os.Create("./doc/ROUTES.md")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.WriteString(content)
		return
	}
	//Otherwise start the Server
	almue.Serve(":8000")
}
