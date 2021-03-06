package almue

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/he4d/almue-backend/model"
)

func (a *Almue) getAllLightingsOfFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor, ok := ctx.Value(floorCtxKey).(*model.Floor)
	if !ok {
		a.logger.Error.Print("Floor from context is not a floor?")
		return
	}

	lightings, err := a.store.GetLightingListOfFloor(floor.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	if err := render.RenderList(w, r, a.newLightingListPayloadResponse(lightings)); err != nil {
		render.Render(w, r, ErrRender(err))
		a.logger.Error.Print(err)
	}
}

func (a *Almue) getAllLightings(w http.ResponseWriter, r *http.Request) {
	lightings, err := a.store.GetLightingList()
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	if err := render.RenderList(w, r, a.newLightingListPayloadResponse(lightings)); err != nil {
		render.Render(w, r, ErrRender(err))
		a.logger.Error.Print(err)
	}
}

func (a *Almue) getLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lighting, ok := ctx.Value(lightingCtxKey).(*model.Lighting)
	if !ok {
		a.logger.Error.Print("Lighting from context is not a lighting?")
		return
	}

	render.Render(w, r, a.newLightingPayloadResponse(lighting))
}

func (a *Almue) createLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor, hasFloorCtx := ctx.Value(floorCtxKey).(*model.Floor)

	l := &lightingPayload{}
	if err := render.Bind(r, l); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		a.logger.Error.Print(err)
	}

	if hasFloorCtx {
		l.FloorID = &floor.ID
	}

	var err error
	l.ID, err = a.store.CreateLighting(l.Lighting)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	lighting, err := a.store.GetLighting(l.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	if err := a.deviceController.RegisterLightings(lighting); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, a.newLightingPayloadResponse(lighting))
}

func (a *Almue) updateLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lighting, ok := ctx.Value(lightingCtxKey).(*model.Lighting)
	if !ok {
		a.logger.Error.Print("Lighting from context is not a lighting?")
		return
	}
	oldLighting := lighting.DeepCopy()

	l := &lightingPayload{Lighting: lighting}
	if err := render.Bind(r, l); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		a.logger.Error.Print(err)
		return
	}

	if l.Lighting.ID != oldLighting.ID {
		err := errors.New("Can not update the lighting to a different id")
		render.Render(w, r, ErrInvalidRequest(err))
		a.logger.Error.Print(err)
		return
	}

	if err := a.store.UpdateLighting(l.Lighting); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	updatedLighting, err := a.store.GetLighting(l.Lighting.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	diffs := oldLighting.GetDifferences(updatedLighting)
	if err := a.deviceController.UpdateLighting(diffs, updatedLighting); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	render.Render(w, r, a.newLightingPayloadResponse(updatedLighting))
}

func (a *Almue) deleteLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lighting, ok := ctx.Value(lightingCtxKey).(*model.Lighting)
	if !ok {
		a.logger.Error.Print("Lighting from context is not a lighting?")
		return
	}

	if err := a.store.DeleteLighting(lighting.ID); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	if err := a.deviceController.UnregisterLighting(lighting.ID); err != nil {
		render.Render(w, r, ErrInternalServer(err))
		a.logger.Error.Print(err)
		return
	}

	render.NoContent(w, r)
}

func (a *Almue) controlLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lighting, ok := ctx.Value(lightingCtxKey).(*model.Lighting)
	if !ok {
		a.logger.Error.Print("Lighting from context is not a lighting?")
		return
	}

	if lighting.Disabled {
		err := errors.New("Device is disabled for controlling")
		render.Render(w, r, ErrInvalidRequest(err))
		a.logger.Info.Print(err)
		return
	}

	action := chi.URLParam(r, "action")
	switch action {
	case "on":
		if err := a.deviceController.TurnLightingOn(lighting.ID); err != nil {
			render.Render(w, r, ErrInternalServer(err))
			a.logger.Error.Print(err)
			return
		}
		break
	case "off":
		if err := a.deviceController.TurnLightingOff(lighting.ID); err != nil {
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
