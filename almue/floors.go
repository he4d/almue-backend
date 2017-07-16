package almue

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
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
	floor, ok := ctx.Value(floorCtxKey).(*model.Floor)
	if !ok {
		a.logger.Error.Print("Floor from context is not a floor?")
		return
	}

	if err := render.Render(w, r, a.newFloorPayloadResponse(floor)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func (a *Almue) updateFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor, ok := ctx.Value(floorCtxKey).(*model.Floor)
	if !ok {
		a.logger.Error.Print("Floor from context is not a floor?")
		return
	}
	oldFloor := floor.DeepCopy()

	f := &floorPayload{Floor: floor}
	if err := render.Bind(r, f); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if f.Floor.ID != oldFloor.ID {
		err := errors.New("Can not update the floor to a different id")
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := a.store.UpdateFloor(f.Floor); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	updatedFloor, err := a.store.GetFloor(f.Floor.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.Render(w, r, a.newFloorPayloadResponse(updatedFloor))
}

func (a *Almue) deleteFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor, ok := ctx.Value(floorCtxKey).(*model.Floor)
	if !ok {
		a.logger.Error.Print("Floor from context is not a floor?")
		return
	}

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
