package almue

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
	"github.com/he4d/almue/embedded"
	"github.com/he4d/almue/model"
)

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
	floor, hasFloorCtx := ctx.Value(floorCtxKey).(*model.Floor)

	s := &shutterPayload{}
	if err := render.Bind(r, s); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
	}

	if hasFloorCtx {
		s.FloorID = &floor.ID
	}

	var err error
	s.ID, err = a.store.CreateShutter(s.Shutter)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	shutter, err := a.store.GetShutter(s.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err := a.deviceController.RegisterShutters(shutter); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	stateSync, err := a.deviceController.GetShutterStateSyncChannels(shutter.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
	}

	go a.startObserveShutterState(shutter.ID, stateSync)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, a.newShutterPayloadResponse(shutter))
}

func (a *Almue) updateShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shutter := ctx.Value(shutterCtxKey).(*model.Shutter)
	oldShutter := shutter.DeepCopy()

	s := &shutterPayload{Shutter: shutter}
	if err := render.Bind(r, s); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if s.Shutter.ID != oldShutter.ID {
		err := errors.New("Can not update the shutter to a different id")
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

	if diffs.HasFlag(model.DIFFDISABLED) {
		if !updatedShutter.Disabled {
			stateSync, err := a.deviceController.GetShutterStateSyncChannels(s.ID)
			if err != nil {
				render.Render(w, r, ErrInternalServer(err))
			}
			go a.startObserveShutterState(updatedShutter.ID, stateSync)
		}
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

	if shutter.Disabled {
		render.Render(w, r, ErrInvalidRequest(errors.New("Device is disabled for controlling")))
		return
	}

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

func (a *Almue) startObserveShutterState(shutterID int64, stateSync *embedded.StateSyncChannels) error {
	for {
		select {
		case newState := <-stateSync.State:
			if err := a.store.UpdateShutterState(shutterID, newState); err != nil {
				//TODO: errorhandling
			}
		case <-stateSync.Quit:
			return nil
		}
	}
}
