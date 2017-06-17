package almue

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/he4d/almue/model"
)

func (a *Almue) getAllFloors(w http.ResponseWriter, r *http.Request) {
	floors, err := a.store.GetFloorList()
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	for i := range floors {
		floors[i].Shutters, err = a.store.GetShutterListOfFloor(floors[i].ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		floors[i].Lightings, err = a.store.GetLightingListOfFloor(floors[i].ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusAccepted, floors)
}

func (a *Almue) createFloor(w http.ResponseWriter, r *http.Request) {
	f := new(model.Floor)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(f); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	_, err := govalidator.ValidateStruct(f)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := a.store.CreateFloor(f)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	floor, err := a.store.GetFloor(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	emptyShutters := make([]*model.Shutter, 0)
	floor.Shutters = emptyShutters

	emptyLightings := make([]*model.Lighting, 0)
	floor.Lightings = emptyLightings

	respondWithJSON(w, http.StatusCreated, floor)
}

func (a *Almue) getFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i, err := strconv.ParseInt(vars["floorID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid floor ID")
		return
	}
	floor, err := a.store.GetFloor(i)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Floor not found")
			return
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	floor.Shutters, err = a.store.GetShutterListOfFloor(i)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	floor.Lightings, err = a.store.GetLightingListOfFloor(i)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, floor)
}

func (a *Almue) updateFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["floorID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid floor ID")
		return
	}

	f := new(model.Floor)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(f); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	f.ID = id

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.store.UpdateFloor(f); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	floor, err := a.store.GetFloor(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, floor)
}

func (a *Almue) deleteFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["floorID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid floor ID")
		return
	}

	//TODO: update a.devices (delete related devices)

	if err := a.store.DeleteFloor(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
