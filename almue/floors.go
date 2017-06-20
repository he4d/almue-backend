package almue

import (
	"encoding/json"
	"net/http"

	"github.com/he4d/almue/model"
)

func (a *Almue) getAllFloors(w http.ResponseWriter, r *http.Request) {
	floors, err := a.store.GetFloorList()
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	for i := range floors {
		floors[i].NumShutters, err = a.store.NumShuttersOfFloor(floors[i].ID)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		floors[i].NumLightings, err = a.store.NumLightingsOfFloor(floors[i].ID)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}

	respondWithJSON(w, http.StatusAccepted, floors)
}

func (a *Almue) createFloor(w http.ResponseWriter, r *http.Request) {
	f := new(model.Floor)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(f); err != nil {
		respondWithError(w, 400)
		return
	}
	defer r.Body.Close()

	id, err := a.store.CreateFloor(f)
	if err != nil {
		respondWithError(w, 500)
		return
	}

	floor, err := a.store.GetFloor(id)
	if err != nil {
		respondWithError(w, 500)
		return
	}

	floor.NumLightings = 0
	floor.NumShutters = 0

	respondWithJSON(w, http.StatusCreated, floor)
}

func (a *Almue) getFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor, ok := ctx.Value(contextKeyFloor).(*model.Floor)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	respondWithJSON(w, http.StatusOK, floor)
}

func (a *Almue) updateFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	oldFloor, ok := ctx.Value(contextKeyFloor).(*model.Floor)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	f := new(model.Floor)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(f); err != nil {
		respondWithError(w, 400)
		return
	}
	defer r.Body.Close()
	f.ID = oldFloor.ID

	if err := a.store.UpdateFloor(f); err != nil {
		respondWithError(w, 500)
		return
	}

	floor, err := a.store.GetFloor(oldFloor.ID)
	if err != nil {
		respondWithError(w, 500)
		return
	}

	respondWithJSON(w, http.StatusOK, floor)
}

func (a *Almue) deleteFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor, ok := ctx.Value(contextKeyFloor).(*model.Floor)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	//TODO: update a.devices (delete related devices)

	if err := a.store.DeleteFloor(floor.ID); err != nil {
		respondWithError(w, 500)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
