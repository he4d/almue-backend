package almue

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/render"
	"github.com/he4d/almue/model"
)

func (a *Almue) getAllFloors(w http.ResponseWriter, r *http.Request) {
	floors, err := a.store.GetFloorList()
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	if err := render.RenderList(w, r, a.newFloorListPayloadResponse(floors)); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

func (a *Almue) createFloor(w http.ResponseWriter, r *http.Request) {
	f := &floorPayload{}
	if err := render.Bind(r, f); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	var err error
	f.ID, err = a.store.CreateFloor(f.Floor)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	floor, err := a.store.GetFloor(f.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, a.newFloorPayloadResponse(floor))
}

func (a *Almue) getFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor := ctx.Value(floorCtxKey).(*model.Floor)

	if err := render.Render(w, r, a.newFloorPayloadResponse(floor)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func (a *Almue) updateFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	oldFloor := ctx.Value(floorCtxKey).(*model.Floor)

	f := new(model.Floor)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(f); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	defer r.Body.Close()
	f.ID = oldFloor.ID

	if err := a.store.UpdateFloor(f); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	f, err := a.store.GetFloor(f.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.Render(w, r, a.newFloorPayloadResponse(f))
}

func (a *Almue) deleteFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor := ctx.Value(floorCtxKey).(*model.Floor)

	shutters, err := a.store.GetShutterListOfFloor(floor.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	for _, shutter := range shutters {
		if err := a.deviceController.UnregisterShutter(shutter.ID); err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}
	}

	lightings, err := a.store.GetLightingListOfFloor(floor.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	for _, lighting := range lightings {
		if err := a.deviceController.UnregisterLighting(lighting.ID); err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}
	}

	if err := a.store.DeleteFloor(floor.ID); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.NoContent(w, r)
}
