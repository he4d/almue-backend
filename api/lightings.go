package api

import (
	"net/http"
)

func (a *API) getAllLightingsOfFloor(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Not implemented")
}

func (a *API) getLighting(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Not implemented")
}

func (a *API) createLighting(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Not implemented")
}

func (a *API) updateLighting(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Not implemented")
}

func (a *API) deleteLighting(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Not implemented")
}
