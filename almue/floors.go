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
	render.Status(r, http.StatusOK)
}

func (a *Almue) createFloor(w http.ResponseWriter, r *http.Request) {
	f := &floorPayload{}
	if err := render.Bind(r, f); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	id, err := a.store.CreateFloor(f.Floor)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	floor := f.Floor
	floor.ID = id

	render.Status(r, http.StatusCreated)
	render.Render(w, r, a.newFloorPayloadResponse(floor)) //TODO: Check err
}

func (a *Almue) getFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor := ctx.Value(floorCtxKey).(*model.Floor)

	if err := render.Render(w, r, a.newFloorPayloadResponse(floor)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

//TODO: hier muss weiter gearbeitet werden!!
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

	render.Render(w, r, a.newFloorPayloadResponse(f)) //TODO: Check err
}

func (a *Almue) deleteFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor := ctx.Value(floorCtxKey).(*model.Floor)

	//TODO: update a.devices (delete related devices)

	if err := a.store.DeleteFloor(floor.ID); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusNoContent)
	render.Render(w, r, a.newNoContentPayloadResponse()) //TODO: Check err
}
