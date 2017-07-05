package almue

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "almue context value " + k.name
}

var (
	floorCtxKey      = &contextKey{"floor"}
	shutterCtxKey    = &contextKey{"shutter"}
	lightingCtxKey   = &contextKey{"lighting"}
	apiVersionCtxKey = &contextKey{"api-version"}
)

func (a *Almue) floorCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		floorID, err := strconv.ParseInt(chi.URLParam(r, "floorID"), 10, 64)
		floor, err := a.store.GetFloor(floorID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			a.logger.Debug.Printf("Failed to put floor to context: %v", err)
			return
		}
		ctx := context.WithValue(r.Context(), floorCtxKey, floor)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Almue) shutterCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shutterID, err := strconv.ParseInt(chi.URLParam(r, "shutterID"), 10, 64)
		shutter, err := a.store.GetShutter(shutterID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			a.logger.Debug.Printf("Failed to put shutter to context: %v", err)
			return
		}
		ctx := context.WithValue(r.Context(), shutterCtxKey, shutter)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Almue) lightingCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lightingID, err := strconv.ParseInt(chi.URLParam(r, "lightingID"), 10, 64)
		lighting, err := a.store.GetLighting(lightingID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			a.logger.Debug.Printf("Failed to put lighting to context: %v", err)
			return
		}
		ctx := context.WithValue(r.Context(), lightingCtxKey, lighting)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func apiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), apiVersionCtxKey, version))
			next.ServeHTTP(w, r)
		})
	}
}
