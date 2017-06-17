package almue

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/he4d/almue/rpi"
	"github.com/he4d/almue/store"
)

// Almue holds all the fields for the complete Application Context
type Almue struct {
	router           *mux.Router
	store            store.Store
	deviceController rpi.DeviceController
	shutterStates    map[int64]*rpi.StateSynchronization
	lightingStates   map[int64]*rpi.StateSynchronization
}

// Initialize sets up the complete api
func (a *Almue) Initialize(dbPath string, simulate bool) {
	a.initializeServer()
	a.initializeDatabase(dbPath)
	a.initializeDeviceController(simulate)
}

// Run must be called to start the api
func (a *Almue) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.router))
}

func (a *Almue) initializeDatabase(dbPath string) {
	a.store = store.New(dbPath)
}

func (a *Almue) initializeDeviceController(simulate bool) {
	a.deviceController = rpi.New(simulate)
	a.shutterStates = make(map[int64]*rpi.StateSynchronization)
	a.lightingStates = make(map[int64]*rpi.StateSynchronization)

	// Register all shutters
	allShutters, err := a.store.GetAllShutters()
	if err != nil {
		log.Println(err)
		log.Fatalln("initializeDeviceController failed")
	}

	shutterStates, err := a.deviceController.RegisterShutters(allShutters...)
	if err != nil {
		log.Println(err)
		log.Fatalln("could not register shutters")
	}

	for id, state := range shutterStates {
		if err := a.registerShutterStateSynchronization(id, state); err != nil {
			log.Println(err)
			log.Fatalln("registering shutterStateSynchronization failed")
		}
	}

	// Register all lightings
	allLightings, err := a.store.GetAllLightings()
	if err != nil {
		log.Println(err)
		log.Fatalln("initializeDeviceController failed")
	}
	lightingStates, err := a.deviceController.RegisterLightings(allLightings...)
	if err != nil {
		log.Println(err)
		log.Fatalln("could not register lightings")
	}

	for id, state := range lightingStates {
		if err := a.registerLightingStateSynchronization(id, state); err != nil {
			log.Println(err)
			log.Fatalln("registering lightingStateSynchronization failed")
		}
	}
}

func (a *Almue) registerShutterStateSynchronization(shutterID int64, stateSync *rpi.StateSynchronization) error {
	if _, ok := a.shutterStates[shutterID]; ok {
		log.Fatalf("shutter with this id: %d is already registered for state synchronization. Exiting...\n", shutterID)
	}
	a.shutterStates[shutterID] = stateSync
	go func() {
		for {
			select {
			case newState := <-stateSync.State:
				if err := a.store.UpdateShutterState(shutterID, newState); err != nil {
					//TODO: errorhandling
				}
			case <-stateSync.Quit:
				return
			}
		}
	}()
	return nil
}

func (a *Almue) registerLightingStateSynchronization(lightingID int64, stateSync *rpi.StateSynchronization) error {
	if _, ok := a.lightingStates[lightingID]; ok {
		log.Fatalf("lighting with this id: %d is already registered for state synchronization. Exiting...\n", lightingID)
	}
	a.lightingStates[lightingID] = stateSync
	go func() {
		for {
			select {
			case newState := <-stateSync.State:
				if err := a.store.UpdateLightingState(lightingID, newState); err != nil {
					//TODO: errorhandling
				}
			case <-stateSync.Quit:
				return
			}
		}
	}()
	return nil
}

func (a *Almue) initializeServer() {
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

	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}/shutters/{shutterID:[0-9]+}/{action:[a-z]+}", a.controlShutter).Methods("POST")
	a.router.HandleFunc("/api/floors/{floorID:[0-9]+}/shutters/{lighting:[0-9]+}/{action:[a-z]+}", a.controlLighting).Methods("POST")
}
