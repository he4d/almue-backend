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
	"github.com/rs/cors"
)

// Almue holds all the fields for the complete Application Context
type Almue struct {
	router           chi.Router
	store            store.Store
	deviceController embedded.DeviceController
	simulate         bool
	dbPath           string
	publicAPI        bool
}

// NewAlmue initializes a new Almue struct, initializes it and return it
func NewAlmue(dbPath string, simulate bool, publicAPI bool) *Almue {
	app := Almue{dbPath: dbPath, simulate: simulate, publicAPI: publicAPI}
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

	// Register all shutters
	allShutters, err := a.store.GetShutterList()
	if err != nil {
		log.Fatalf("initializeDeviceController failed %s", err)
	}

	if err := a.deviceController.RegisterShutters(allShutters...); err != nil {
		log.Fatalf("could not register shutters %s", err)
	}

	for _, shutter := range allShutters {
		syncState, err := a.deviceController.GetShutterStateSyncChannels(shutter.ID)
		if err != nil {
			log.Fatalf("could not get shutterStateSync: %s", err)
		}
		go a.startObserveShutterState(shutter.ID, syncState)
	}

	// Register all lightings
	allLightings, err := a.store.GetLightingList()
	if err != nil {
		log.Fatalf("initializeDeviceController failed %s", err)
	}
	if err := a.deviceController.RegisterLightings(allLightings...); err != nil {
		log.Fatalf("could not register lightings %s", err)
	}

	for _, lighting := range allLightings {
		syncState, err := a.deviceController.GetLightingStateSyncChannels(lighting.ID)
		if err != nil {
			log.Fatalf("could not get lightingStateSync: %s", err)
		}
		go a.startObserveLightingState(lighting.ID, syncState)
	}
}

func (a *Almue) initializeRouter() {
	a.router = chi.NewRouter()

	// Set up the middleware
	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)

	if a.publicAPI {
		cors := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		})
		a.router.Use(cors.Handler)
	}
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
				r.Get("/", a.getAllLightings)
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
