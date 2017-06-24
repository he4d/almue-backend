package almue

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
	"github.com/he4d/almue/model"
)

func (a *Almue) getAllLightingsOfFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor := ctx.Value(floorCtxKey).(*model.Floor)

	lightings, err := a.store.GetLightingListOfFloor(floor.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err := render.RenderList(w, r, a.newLightingListPayloadResponse(lightings)); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

func (a *Almue) getAllLightings(w http.ResponseWriter, r *http.Request) {
	lightings, err := a.store.GetLightingList()
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err := render.RenderList(w, r, a.newLightingListPayloadResponse(lightings)); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

func (a *Almue) getLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lighting := ctx.Value(lightingCtxKey).(*model.Lighting)

	render.Render(w, r, a.newLightingPayloadResponse(lighting))
}

func (a *Almue) createLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor := ctx.Value(floorCtxKey).(*model.Floor)

	l := new(model.Lighting)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(l); err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}
	defer r.Body.Close()
	l.FloorID = floor.ID

	var err error
	l.ID, err = a.store.CreateLighting(l)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	l, err = a.store.GetLighting(l.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err := a.deviceController.RegisterLightings(l); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	syncState, err := a.deviceController.GetLightingStateSyncChannels(l.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	go a.startObserveLightingState(l.ID, syncState)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, a.newLightingPayloadResponse(l))
}

func (a *Almue) updateLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lighting := ctx.Value(lightingCtxKey).(*model.Lighting)
	oldLighting := lighting.DeepCopy()

	l := &lightingPayload{Lighting: lighting}
	if err := render.Bind(r, l); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := a.store.UpdateLighting(l.Lighting); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	updatedLighting, err := a.store.GetLighting(l.Lighting.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	diffs := oldLighting.GetDifferences(updatedLighting)
	if err := a.deviceController.UpdateLighting(diffs, updatedLighting); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.Render(w, r, a.newLightingPayloadResponse(updatedLighting))
}

func (a *Almue) deleteLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lighting := ctx.Value(lightingCtxKey).(*model.Lighting)

	if err := a.store.DeleteLighting(lighting.ID); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err := a.deviceController.UnregisterLighting(lighting.ID); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.NoContent(w, r)
}

func (a *Almue) controlLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lighting := ctx.Value(shutterCtxKey).(*model.Lighting)

	action := chi.URLParam(r, "action")
	switch action {
	case "on":
		if err := a.deviceController.TurnLightingOn(lighting.ID); err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}
		break
	case "off":
		if err := a.deviceController.TurnLightingOff(lighting.ID); err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}
		break
	default:
		render.Render(w, r, ErrInvalidRequest(errors.New("Action not supported")))
		return
	}
	render.NoContent(w, r)
}
