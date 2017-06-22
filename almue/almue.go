package almue

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/docgen"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/render"
	"github.com/he4d/almue/embedded"
	"github.com/he4d/almue/store"
)

// Almue holds all the fields for the complete Application Context
type Almue struct {
	router           chi.Router
	store            store.Store
	deviceController embedded.DeviceController
	shutterStates    map[int64]*embedded.StateSynchronization
	lightingStates   map[int64]*embedded.StateSynchronization
	simulate         bool
	dbPath           string
}

// NewAlmue initializes a new Almue struct, initializes it and return it
func NewAlmue(dbPath string, simulate bool) *Almue {
	app := Almue{dbPath: dbPath, simulate: simulate}
	app.initialize()
	return &app
}

// Serve must be called to start the Almue backend
func (a *Almue) Serve(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.router))
}

//GenerateRoutesDoc generates a markdown documentation of the API routes
func (a *Almue) GenerateRoutesDoc() {
	content := (docgen.MarkdownRoutesDoc(a.router, docgen.MarkdownOpts{
		ProjectPath: "github.com/he4d/almue",
		Intro:       "Welcome to the Almue generated docs.",
	}))
	f, err := os.Create("./doc/ROUTES.md")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(content)
}

func (a *Almue) initialize() {
	a.initializeRouter()
	a.initializeDatabase()
	a.initializeDeviceController()
}

func (a *Almue) initializeDatabase() {
	a.store = store.New(a.dbPath)
}

func (a *Almue) initializeDeviceController() {
	a.deviceController = embedded.New(a.simulate)
	a.shutterStates = make(map[int64]*embedded.StateSynchronization)
	a.lightingStates = make(map[int64]*embedded.StateSynchronization)

	// Register all shutters
	allShutters, err := a.store.GetShutterList()
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
	allLightings, err := a.store.GetLightingList()
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

func (a *Almue) registerShutterStateSynchronization(shutterID int64, stateSync *embedded.StateSynchronization) error {
	if _, ok := a.shutterStates[shutterID]; ok {
		log.Fatalf("shutter with this id: %d is already registered for state synchronization.", shutterID)
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

func (a *Almue) registerLightingStateSynchronization(lightingID int64, stateSync *embedded.StateSynchronization) error {
	if _, ok := a.lightingStates[lightingID]; ok {
		log.Fatalf("lighting with this id: %d is already registered for state synchronization.", lightingID)
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

func (a *Almue) initializeRouter() {
	a.router = chi.NewRouter()

	// Set up the middleware
	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)
	a.router.Use(render.SetContentType(render.ContentTypeJSON))

	// Serve static files
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "frontend/dist")
	fileServer(a.router, "/", http.Dir(filesDir))

	// API version 1
	a.router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Use(apiVersionCtx("v1"))
			r.Route("/shutters", func(r chi.Router) {
				r.Get("/", a.getAllShutters)
			})
			r.Route("/lightings", func(r chi.Router) {
				r.Get("/", a.getAllLightings)
			})
			r.Route("/floors", func(r chi.Router) {
				r.Get("/", a.getAllFloors)
				r.Post("/", a.createFloor)
				r.Route("/{floorID:[0-9]+$}", func(r chi.Router) {
					r.Use(a.floorCtx)
					r.Get("/", a.getFloor)
					r.Put("/", a.updateFloor)
					r.Delete("/", a.deleteFloor)
					r.Route("/shutters", func(r chi.Router) {
						r.Get("/", a.getAllShuttersOfFloor)
						r.Post("/", a.createShutter)
						r.Route("/{shutterID:[0-9]+$}", func(r chi.Router) {
							r.Use(a.shutterCtx)
							r.Get("/", a.getShutter)
							r.Put("/", a.updateShutter)
							r.Delete("/", a.deleteShutter)
							r.Route("/{action:[a-z]+$}", func(r chi.Router) {
								r.Post("/", a.controlShutter)
							})
						})
					})
					r.Route("/lightings", func(r chi.Router) {
						r.Get("/", a.getAllLightingsOfFloor)
						r.Post("/", a.createLighting)
						r.Route("/{lightingID:[0-9]+$}", func(r chi.Router) {
							r.Use(a.lightingCtx)
							r.Get("/", a.getLighting)
							r.Put("/", a.updateLighting)
							r.Delete("/", a.deleteLighting)
							r.Route("/{action:[a-z]+$}", func(r chi.Router) {
								r.Post("/", a.controlLighting)
							})
						})
					})
				})
			})
		})
	})
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, ":*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
