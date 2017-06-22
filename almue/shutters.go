package almue

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
	"github.com/he4d/almue/model"
)

func (a *Almue) newShutterPayloadResponse(shutter *model.Shutter) *shutterPayload {
	resp := &shutterPayload{Shutter: shutter}

	return resp
}

func (a *Almue) getAllShuttersOfFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor := ctx.Value(floorCtxKey).(*model.Floor)

	shutters, err := a.store.GetShutterListOfFloor(floor.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err := render.RenderList(w, r, a.newShutterListPayloadResponse(shutters)); err != nil {
		render.Render(w, r, ErrRender(err))
	}
	render.Status(r, http.StatusOK)
}

func (a *Almue) getAllShutters(w http.ResponseWriter, r *http.Request) {
	shutters, err := a.store.GetShutterList()
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err := render.RenderList(w, r, a.newShutterListPayloadResponse(shutters)); err != nil {
		render.Render(w, r, ErrRender(err))
	}
	render.Status(r, http.StatusOK)
}

func (a *Almue) getShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shutter := ctx.Value(shutterCtxKey).(*model.Shutter)

	render.Status(r, http.StatusOK)
	render.Render(w, r, a.newShutterPayloadResponse(shutter)) //TODO: Check err
}

func (a *Almue) createShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor := ctx.Value(floorCtxKey).(*model.Floor)

	s := new(model.Shutter)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(s); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	defer r.Body.Close()
	s.FloorID = floor.ID

	newID, err := a.store.CreateShutter(s)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	s.ID = newID

	shutterState, err := a.deviceController.RegisterShutters(s)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err = a.registerShutterStateSynchronization(newID, shutterState[newID]); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, a.newShutterPayloadResponse(s)) //TODO: Check err
}

func (a *Almue) updateShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shutter := ctx.Value(shutterCtxKey).(*model.Shutter)

	s := &shutterPayload{Shutter: shutter}
	if err := render.Bind(r, s); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	shutter = s.Shutter

	if err := a.store.UpdateShutter(shutter); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	//TODO: update a.devices

	render.Status(r, http.StatusOK)
	render.Render(w, r, a.newShutterPayloadResponse(shutter)) //TODO: Check err
}

func (a *Almue) deleteShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shutter, ok := ctx.Value(shutterCtxKey).(*model.Shutter)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	if err := a.store.DeleteShutter(shutter.ID); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err := a.deviceController.UnregisterShutter(shutter.ID); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusNoContent)
	render.Render(w, r, a.newNoContentPayloadResponse()) //TODO: Check err
}

func (a *Almue) controlShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shutter := ctx.Value(shutterCtxKey).(*model.Shutter)

	action := chi.URLParam(r, "action")
	switch action {
	case "open":
		if err := a.deviceController.OpenShutter(shutter.ID); err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}
		break
	case "close":
		if err := a.deviceController.CloseShutter(shutter.ID); err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}
		break
	case "stop":
		if err := a.deviceController.StopShutter(shutter.ID); err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}
		break
	default:
		render.Render(w, r, ErrInvalidRequest(errors.New("Action not supported")))
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, a.newNoContentPayloadResponse())
}
