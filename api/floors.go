package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *API) getAllFloors(w http.ResponseWriter, r *http.Request) {
	floors, err := a.store.GetFloorList()
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	}
	respondWithJSON(w, http.StatusAccepted, floors)
}

func (a *API) createFloor(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Not implemented")
}

func (a *API) getFloor(w http.ResponseWriter, r *http.Request) {
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
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, floor)
}

func (a *API) updateFloor(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Not implemented")
}

func (a *API) deleteFloor(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Not implemented")
}
