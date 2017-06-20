package almue

import (
	"context"
	"net/http"
	"strconv"

	"github.com/pressly/chi"
)

type contextKey string

var (
	contextKeyFloor        = contextKey("floor")
	contextKeyShutter      = contextKey("shutter")
	contextKeyLighting     = contextKey("lighting")
	contextKeyDeviceAction = contextKey("device-action")
)

func (a *Almue) floorCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		floorID, err := strconv.ParseInt(chi.URLParam(r, "floorID"), 10, 64)
		floor, err := a.store.GetFloor(floorID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		numShutters, err := a.store.NumShuttersOfFloor(floorID)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		floor.NumShutters = numShutters
		numLightings, err := a.store.NumLightingsOfFloor(floorID)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		floor.NumLightings = numLightings
		ctx := context.WithValue(r.Context(), contextKeyFloor, floor)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Almue) deviceActionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		action := chi.URLParam(r, "action")
		ctx := context.WithValue(r.Context(), contextKeyDeviceAction, action)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Almue) shutterCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shutterID, err := strconv.ParseInt(chi.URLParam(r, "shutterID"), 10, 64)
		shutter, err := a.store.GetShutter(shutterID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyShutter, shutter)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Almue) lightingCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lightingID, err := strconv.ParseInt(chi.URLParam(r, "lightingID"), 10, 64)
		lighting, err := a.store.GetLighting(lightingID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyLighting, lighting)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
