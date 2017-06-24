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
}

func (a *Almue) getShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shutter := ctx.Value(shutterCtxKey).(*model.Shutter)

	render.Render(w, r, a.newShutterPayloadResponse(shutter))
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

	var err error
	s.ID, err = a.store.CreateShutter(s)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	s, err = a.store.GetShutter(s.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err := a.deviceController.RegisterShutters(s); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	stateSync, err := a.deviceController.GetShutterStateSyncChannels(s.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
	}

	go a.startObserveShutterState(s.ID, stateSync)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, a.newShutterPayloadResponse(s))
}

func (a *Almue) updateShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	oldShutter := ctx.Value(shutterCtxKey).(*model.Shutter)

	s := &shutterPayload{Shutter: oldShutter}
	if err := render.Bind(r, s); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := a.store.UpdateShutter(s.Shutter); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	updatedShutter, err := a.store.GetShutter(s.Shutter.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	diffs := oldShutter.GetDifferences(updatedShutter)
	if err := a.deviceController.UpdateShutter(diffs, updatedShutter); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.Render(w, r, a.newShutterPayloadResponse(updatedShutter))
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

	render.NoContent(w, r)
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
	render.NoContent(w, r)
}
