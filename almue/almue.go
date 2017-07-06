package almue

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/docgen"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/render"
	"github.com/he4d/simplejack"
	"github.com/rs/cors"
)

// Almue holds all the fields for the complete Application Context
type Almue struct {
	router           chi.Router
	server           *http.Server
	store            DeviceStore
	deviceController DeviceController
	simulate         bool
	publicAPI        bool
	logger           *simplejack.Logger
}

// New initializes a new Almue struct, initializes it and return it
func New(store DeviceStore, deviceController DeviceController, logger *simplejack.Logger, publicAPI bool) *Almue {
	app := &Almue{store: store, deviceController: deviceController, logger: logger, publicAPI: publicAPI}
	logger.Info.Print("Initializing the application")
	if err := app.initialize(); err != nil {
		logger.Fatal.Fatal(err)
	}
	logger.Info.Print("Successfully initialized the application")
	return app
}

// Serve must be called to start the Almue backend
// if an error occured on calling listenAndServe the returned
// error chan will contain the error message
func (a *Almue) Serve(addr string) <-chan error {
	serveErrorChan := make(chan error)
	go func() {
		a.server = &http.Server{Addr: addr, Handler: a.router}
		if err := a.server.ListenAndServe(); err != nil {
			serveErrorChan <- err
		}
	}()
	return serveErrorChan
}

func (a *Almue) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
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

func (a *Almue) initialize() error {
	if err := a.initializeRouter(); err != nil {
		return err
	}
	if err := a.feedDeviceController(); err != nil {
		return err
	}
	return nil
}

func (a *Almue) feedDeviceController() error {
	allShutters, err := a.store.GetShutterList()
	if err != nil {
		return err
	}

	if err := a.deviceController.RegisterShutters(allShutters...); err != nil {
		return err
	}

	allLightings, err := a.store.GetLightingList()
	if err != nil {
		return err
	}

	if err := a.deviceController.RegisterLightings(allLightings...); err != nil {
		return err
	}
	return nil
}

func (a *Almue) initializeRouter() error {
	a.router = chi.NewRouter()

	// Set up the middleware
	//TODO: ONLY FOR DEBUGGING
	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.Logger)
	//
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
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}
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
	return nil
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
