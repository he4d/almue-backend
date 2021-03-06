package almue

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/he4d/almue-backend/model"
)

func (a *Almue) getAllShuttersOfFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor, ok := ctx.Value(floorCtxKey).(*model.Floor)
	if !ok {
		a.logger.Error.Print("Floor from context is not a floor?")
		return
	}

	shutters, err := a.store.GetShutterListOfFloor(floor.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	if err := render.RenderList(w, r, a.newShutterListPayloadResponse(shutters)); err != nil {
		render.Render(w, r, ErrRender(err))
		a.logger.Error.Print(err)
	}
}

func (a *Almue) getAllShutters(w http.ResponseWriter, r *http.Request) {
	shutters, err := a.store.GetShutterList()
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	if err := render.RenderList(w, r, a.newShutterListPayloadResponse(shutters)); err != nil {
		render.Render(w, r, ErrRender(err))
		a.logger.Error.Print(err)
	}
}

func (a *Almue) getShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shutter, ok := ctx.Value(shutterCtxKey).(*model.Shutter)
	if !ok {
		a.logger.Error.Print("Shutter from context is not a shutter?")
		return
	}

	render.Render(w, r, a.newShutterPayloadResponse(shutter))
}

func (a *Almue) createShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor, hasFloorCtx := ctx.Value(floorCtxKey).(*model.Floor)

	s := &shutterPayload{}
	if err := render.Bind(r, s); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		a.logger.Error.Print(err)
	}

	if hasFloorCtx {
		s.FloorID = &floor.ID
	}

	var err error
	s.ID, err = a.store.CreateShutter(s.Shutter)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	shutter, err := a.store.GetShutter(s.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	if err := a.deviceController.RegisterShutters(shutter); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, a.newShutterPayloadResponse(shutter))
}

func (a *Almue) updateShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shutter, ok := ctx.Value(shutterCtxKey).(*model.Shutter)
	if !ok {
		a.logger.Error.Print("Shutter from context is not a shutter?")
		return
	}

	oldShutter := shutter.DeepCopy()

	s := &shutterPayload{Shutter: shutter}
	if err := render.Bind(r, s); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		a.logger.Error.Print(err)
		return
	}

	if s.Shutter.ID != oldShutter.ID {
		err := errors.New("Can not update the shutter to a different id")
		render.Render(w, r, ErrInvalidRequest(err))
		a.logger.Error.Print(err)
		return
	}

	if err := a.store.UpdateShutter(s.Shutter); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	updatedShutter, err := a.store.GetShutter(s.Shutter.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	diffs := oldShutter.GetDifferences(updatedShutter)
	if err := a.deviceController.UpdateShutter(diffs, updatedShutter); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	render.Render(w, r, a.newShutterPayloadResponse(updatedShutter))
}

func (a *Almue) deleteShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shutter, ok := ctx.Value(shutterCtxKey).(*model.Shutter)
	if !ok {
		a.logger.Error.Print("Shutter from context is not a shutter?")
		return
	}

	if err := a.store.DeleteShutter(shutter.ID); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	if err := a.deviceController.UnregisterShutter(shutter.ID); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	render.NoContent(w, r)
}

func (a *Almue) controlShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shutter, ok := ctx.Value(shutterCtxKey).(*model.Shutter)
	if !ok {
		a.logger.Error.Print("Shutter from context is not a shutter?")
		return
	}

	if shutter.Disabled {
		err := errors.New("Device is disabled for controlling")
		render.Render(w, r, ErrInvalidRequest(err))
		a.logger.Info.Print(err)
		return
	}

	action := chi.URLParam(r, "action")
	switch action {
	case "open":
		if err := a.deviceController.OpenShutter(shutter.ID); err != nil {
			render.Render(w, r, ErrInternalServer(err))
			a.logger.Error.Print(err)
			return
		}
		break
	case "close":
		if err := a.deviceController.CloseShutter(shutter.ID); err != nil {
			render.Render(w, r, ErrInternalServer(err))
			a.logger.Error.Print(err)
			return
		}
		break
	case "stop":
		if err := a.deviceController.StopShutter(shutter.ID); err != nil {
			render.Render(w, r, ErrInternalServer(err))
			a.logger.Error.Print(err)
			return
		}
		break
	default:
		err := errors.New("Action not supported")
		render.Render(w, r, ErrInvalidRequest(err))
		a.logger.Info.Print(err)
		return
	}
	render.NoContent(w, r)
}
