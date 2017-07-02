package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/he4d/almue/almue"
	"github.com/he4d/almue/store"
	"github.com/he4d/simplejack"
	_ "github.com/mattn/go-sqlite3"
)

var (
	simulate    = flag.Bool("simulate", false, "starts simulation mode without gpio (operations will be written to stdout instead)")
	routes      = flag.Bool("routes", false, "generate router documentation")
	publicAPI   = flag.Bool("publicapi", false, "enables public access to the rest service")
	logLevel    = flag.Int("loglevel", 3, "set the minimum loglevel 0 = Trace, 1 = Debug, 2 = Info, 3 = Warning, 4 = Error, 5 = Fatal")
	logToStdout = flag.Bool("logtostdout", false, "set this to true to get logging to the stdout instead of a logfile")
)

func main() {
	flag.Parse()
	if *logLevel < 0 || *logLevel > 5 {
		log.Fatalf("Log level must be between 0 and 5!")
	}
	sjLogLevel := simplejack.LogLevel(*logLevel)

	var writer io.Writer
	if *logToStdout {
		writer = os.Stdout
	} else {
		var err error
		writer, err = os.OpenFile("almue.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file almue.log : %s", err)
		}
	}

	logger := simplejack.New(sjLogLevel, writer)

	store, err := store.New("./almue.db", logger)
	if err != nil {
		logger.Fatal.Fatalf("Could not create a new store: %v", err)
	}

	//TODO: Create the deviceController here and inject to the new almue instance

	almue := almue.New(store, logger, *simulate, *publicAPI)
	if *routes {
		almue.GenerateRoutesDoc()
		return
	}
	almue.Serve(":8000")
}
