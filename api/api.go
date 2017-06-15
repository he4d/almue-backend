package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/he4d/almue/store"
)

// API holds all the fields for the complete API Context
type API struct {
	router *mux.Router
	store  store.Store
}

// Initialize sets up the complete api
func (a *API) Initialize(dbPath string) {
	a.initialize()
	a.initializeDatabase(dbPath)
}

// Run must be called to start the api
func (a *API) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.router))
}

func (a *API) initializeDatabase(dbPath string) {
	a.store = store.New(dbPath)
}

func (a *API) initializeStaticContent() {
}

func (a *API) initialize() {
	a.router = mux.NewRouter()

	// Serve static files
	fs := http.FileServer(http.Dir("frontend/static"))
	a.router.Handle("/", fs)

	// Serve the api functions
	a.router.HandleFunc("/api/floors", a.getAllFloors).Methods("GET")
	a.router.HandleFunc("/api/floors", a.createFloor).Methods("POST")
	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}", a.getFloor).Methods("GET")
	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}", a.updateFloor).Methods("PUT")
	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}", a.deleteFloor).Methods("DELETE")
	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}/shutters", a.getAllShuttersOfFloor).Methods("GET")
	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}/shutters", a.createShutter).Methods("POST")
	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}/shutters/{shutterID:[0-9]+}", a.getShutter).Methods("GET")
	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}/shutters/{shutterID:[0-9]+}", a.updateShutter).Methods("PUT")
	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}/shutters/{shutterID:[0-9]+}", a.deleteShutter).Methods("DELETE")
	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}/lightings", a.getAllLightingsOfFloor).Methods("GET")
	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}/lightings", a.createLighting).Methods("POST")
	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}/lightings/{lightingID:[0-9]+}", a.getLighting).Methods("GET")
	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}/lightings/{lightingID:[0-9]+}", a.updateLighting).Methods("PUT")
	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}/lightings/{lightingID:[0-9]+}", a.deleteLighting).Methods("DELETE")
}
