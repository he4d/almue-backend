package api

import (
	"net/http"
)

func (a *API) getAllShuttersOfFloor(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Not implemented")
}

func (a *API) getShutter(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Not implemented")
}

func (a *API) createShutter(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Not implemented")
}

func (a *API) updateShutter(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Not implemented")
}

func (a *API) deleteShutter(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Not implemented")
}
