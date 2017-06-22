package almue

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
	"github.com/he4d/almue/model"
)

type shutterPayload struct {
	*model.Shutter
}

type shutterListPayload []*shutterPayload

func (s *shutterPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *shutterPayload) Bind(r *http.Request) error {
	return nil
}

func (a *Almue) newShutterListPayloadResponse(shutters []*model.Shutter) []render.Renderer {
	list := []render.Renderer{}
	for _, shutter := range shutters {
		list = append(list, a.newShutterPayloadResponse(shutter))
	}
	return list
}

func (a *Almue) newShutterPayloadResponse(shutter *model.Shutter) *shutterPayload {
	resp := &shutterPayload{Shutter: shutter}

	return resp
}

func (a *Almue) getAllShuttersOfFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor, ok := ctx.Value(contextKeyFloor).(*model.Floor)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	shutters, err := a.store.GetShutterListOfFloor(floor.ID)
	if err != nil {
		respondWithError(w, 500)
		return
	}
	respondWithJSON(w, http.StatusOK, shutters)
}

func (a *Almue) getAllShutters(w http.ResponseWriter, r *http.Request) {
	shutters, err := a.store.GetShutterList()
	if err != nil {
		respondWithError(w, 500)
		return
	}
	respondWithJSON(w, http.StatusOK, shutters)
}

func (a *Almue) getShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shutter, ok := ctx.Value(contextKeyShutter).(*model.Shutter)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	respondWithJSON(w, http.StatusOK, shutter)
}

func (a *Almue) createShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor, ok := ctx.Value(contextKeyFloor).(*model.Floor)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	s := new(model.Shutter)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(s); err != nil {
		respondWithError(w, 400)
		return
	}
	defer r.Body.Close()
	s.FloorID = floor.ID

	newID, err := a.store.CreateShutter(s)
	if err != nil {
		respondWithError(w, 500)
		return
	}

	shutter, err := a.store.GetShutter(newID)
	if err != nil {
		respondWithError(w, 500)
		return
	}

	shutterState, err := a.deviceController.RegisterShutters(shutter)
	if err != nil {
		respondWithError(w, 500)
		return
	}

	if err = a.registerShutterStateSynchronization(newID, shutterState[newID]); err != nil {
		respondWithError(w, 500)
		return
	}

	respondWithJSON(w, http.StatusCreated, shutter)
}

func (a *Almue) updateShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shutter, ok := ctx.Value(contextKeyShutter).(*model.Shutter)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	s := new(model.Shutter)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(s); err != nil {
		respondWithError(w, 400)
		return
	}
	defer r.Body.Close()

	if err := a.store.UpdateShutter(s); err != nil {
		respondWithError(w, 500)
		return
	}

	shutter, err := a.store.GetShutter(shutter.ID)
	if err != nil {
		respondWithError(w, 500)
		return
	}

	//TODO: update a.devices

	respondWithJSON(w, http.StatusOK, shutter)
}

func (a *Almue) deleteShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shutter, ok := ctx.Value(contextKeyShutter).(*model.Shutter)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	if err := a.store.DeleteShutter(shutter.ID); err != nil {
		respondWithError(w, 500)
		return
	}

	if err := a.deviceController.UnregisterShutter(shutter.ID); err != nil {
		respondWithError(w, 500)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (a *Almue) controlShutter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shutter, ok := ctx.Value(contextKeyShutter).(*model.Shutter)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	_, err := a.store.GetShutter(shutter.ID)
	if err != nil {
		respondWithError(w, 404)
		return
	}

	action := chi.URLParam(r, "action")
	switch action {
	case "open":
		if err := a.deviceController.OpenShutter(shutter.ID); err != nil {
			respondWithError(w, 400)
			return
		}
		break
	case "close":
		if err := a.deviceController.CloseShutter(shutter.ID); err != nil {
			respondWithError(w, 400)
			return
		}
		break
	case "stop":
		if err := a.deviceController.StopShutter(shutter.ID); err != nil {
			respondWithError(w, 400)
			return
		}
		break
	default:
		respondWithError(w, 400)
		return
	}
}
