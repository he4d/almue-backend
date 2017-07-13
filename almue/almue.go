package almue

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"time"

	"io/ioutil"

	"bytes"

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
func New(store DeviceStore, deviceController DeviceController, logger *simplejack.Logger, publicAPI bool) (*Almue, error) {
	app := &Almue{store: store, deviceController: deviceController, logger: logger, publicAPI: publicAPI}
	if err := app.initialize(); err != nil {
		return nil, err
	}
	return app, nil
}

// Serve must be called to start the Almue backend
// if an error occured on calling ListenAndServe the returned
// error chan will contain the error message
func (a *Almue) Serve(addr string) <-chan error {
	serveError := make(chan error)
	go func() {
		a.server = &http.Server{Addr: addr, Handler: a.router}
		if err := a.server.ListenAndServe(); err != nil {
			serveError <- err
		}
	}()
	return serveError
}

func (a *Almue) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error.Printf("Could not shutdown the server: %v", err)
		return
	}
	a.logger.Info.Print("almue server stopped successfully")
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
			r.Route("/manage", func(r chi.Router) {
				r.Get("/logfile", a.getLogfile)
				r.Route("/db", func(r chi.Router) {
					r.Get("/backup", a.retrieveStoreBackup)
					r.Post("/backup", a.restoreStoreBackup)
				})
			})
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

func (a *Almue) getLogfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	file, err := ioutil.ReadFile("almue.log")
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	w.Write(file)
}

func (a *Almue) retrieveStoreBackup(w http.ResponseWriter, r *http.Request) {
	file, err := a.store.GetBackup()
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\"almue.db\"")
	w.WriteHeader(http.StatusOK)
	w.Write(file)
}

func (a *Almue) restoreStoreBackup(w http.ResponseWriter, r *http.Request) {
	readForm, err := r.MultipartReader()
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
	}
	for {
		part, err := readForm.NextPart()
		if err == io.EOF {
			break
		}
		if part.FormName() == "dbfile" {
			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(part); err != nil {
				render.Render(w, r, ErrInternalServer(err))
			}
			if err := a.store.RestoreBackup(buf.Bytes()); err != nil {
				render.Render(w, r, ErrInternalServer(err))
			}
			render.NoContent(w, r)
		}
	}
	render.Render(w, r, ErrInvalidRequest(err))
}
