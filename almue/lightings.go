package almue

import (
	"encoding/json"
	"net/http"

	"github.com/he4d/almue/model"
	"github.com/pressly/chi"
)

func (a *Almue) getAllLightingsOfFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor, ok := ctx.Value(contextKeyFloor).(*model.Floor)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	lightings, err := a.store.GetLightingListOfFloor(floor.ID)
	if err != nil {
		respondWithError(w, 500)
		return
	}
	respondWithJSON(w, http.StatusOK, lightings)
}

func (a *Almue) getAllLightings(w http.ResponseWriter, r *http.Request) {
	lightings, err := a.store.GetLightingList()
	if err != nil {
		respondWithError(w, 500)
		return
	}
	respondWithJSON(w, http.StatusOK, lightings)
}

func (a *Almue) getLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lighting, ok := ctx.Value(contextKeyLighting).(*model.Lighting)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	respondWithJSON(w, http.StatusOK, lighting)
}

func (a *Almue) createLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor, ok := ctx.Value(contextKeyFloor).(*model.Floor)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	l := new(model.Lighting)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(l); err != nil {
		respondWithError(w, 400)
		return
	}
	defer r.Body.Close()
	l.FloorID = floor.ID

	newID, err := a.store.CreateLighting(l)
	if err != nil {
		respondWithError(w, 500)
		return
	}

	lighting, err := a.store.GetLighting(newID)
	if err != nil {
		respondWithError(w, 500)
		return
	}

	lightingState, err := a.deviceController.RegisterLightings(lighting)
	if err != nil {
		respondWithError(w, 500)
		return
	}

	if err = a.registerLightingStateSynchronization(newID, lightingState[newID]); err != nil {
		respondWithError(w, 500)
		return
	}

	respondWithJSON(w, http.StatusCreated, lighting)
}

func (a *Almue) updateLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lighting, ok := ctx.Value(contextKeyLighting).(*model.Lighting)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	l := new(model.Lighting)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(l); err != nil {
		respondWithError(w, 400)
		return
	}
	defer r.Body.Close()

	if err := a.store.UpdateLighting(l); err != nil {
		respondWithError(w, 500)
		return
	}

	lighting, err := a.store.GetLighting(lighting.ID)
	if err != nil {
		respondWithError(w, 500)
		return
	}

	//TODO: update a.devices

	respondWithJSON(w, http.StatusOK, lighting)
}

func (a *Almue) deleteLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lighting, ok := ctx.Value(contextKeyLighting).(*model.Lighting)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	if err := a.store.DeleteLighting(lighting.ID); err != nil {
		respondWithError(w, 500)
		return
	}

	if err := a.deviceController.UnregisterLighting(lighting.ID); err != nil {
		respondWithError(w, 500)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (a *Almue) controlLighting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lighting, ok := ctx.Value(contextKeyShutter).(*model.Lighting)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	_, err := a.store.GetLighting(lighting.ID)
	if err != nil {
		respondWithError(w, 404)
		return
	}

	action := chi.URLParam(r, "action")
	switch action {
	case "on":
		if err := a.deviceController.TurnLightingOn(lighting.ID); err != nil {
			respondWithError(w, 400)
			return
		}
		break
	case "off":
		if err := a.deviceController.TurnLightingOff(lighting.ID); err != nil {
			respondWithError(w, 400)
			return
		}
		break
	default:
		respondWithError(w, 400)
		return
	}
}
