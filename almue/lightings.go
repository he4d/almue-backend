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

func (a *Almue) getAllLightingsOfFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i, err := strconv.ParseInt(vars["floorID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid floor ID")
		return
	}
	lightings, err := a.store.GetLightingListOfFloor(i)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, lightings)
}

func (a *Almue) getLighting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	floorID, err := strconv.ParseInt(vars["floorID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid floor ID")
		return
	}
	shutterID, err := strconv.ParseInt(vars["lightingID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid lighting ID")
		return
	}
	lighting, err := a.store.GetLightingByFloor(shutterID, floorID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Lighting not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, lighting)
}

func (a *Almue) createLighting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i, err := strconv.ParseInt(vars["floorID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid floor ID")
		return
	}
	l := new(model.Lighting)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(l); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	l.FloorID = i

	_, err = govalidator.ValidateStruct(l)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	newID, err := a.store.CreateLighting(l)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	lighting, err := a.store.GetLightingByFloor(newID, i)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	if err := a.deviceController.RegisterLightings(lighting); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusCreated, lighting)
}

func (a *Almue) updateLighting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	floorID, err := strconv.ParseInt(vars["floorID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid floor ID")
		return
	}

	lightingID, err := strconv.ParseInt(vars["lightingID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid lighting ID")
		return
	}

	l := new(model.Lighting)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(l); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()

	_, err = govalidator.ValidateStruct(l)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.store.UpdateLighting(l); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	lighting, err := a.store.GetLightingByFloor(lightingID, floorID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	//TODO: update a.devices

	respondWithJSON(w, http.StatusOK, lighting)
}

func (a *Almue) deleteLighting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lightingID, err := strconv.ParseInt(vars["lightingID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid lighting ID")
		return
	}

	if err := a.store.DeleteLighting(lightingID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := a.deviceController.UnregisterLighting(lightingID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (a *Almue) controlLighting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	floorID, err := strconv.ParseInt(vars["floorID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid floor ID")
		return
	}
	lightingID, err := strconv.ParseInt(vars["lightingID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid lighting ID")
		return
	}

	_, err = a.store.GetLightingByFloor(lightingID, floorID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Shutter not found")
		return
	}

	action := vars["action"]
	switch action {
	case "on":
		if err := a.deviceController.TurnLightingOn(lightingID); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		break
	case "off":
		if err := a.deviceController.TurnLightingOff(lightingID); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		break
	default:
		respondWithError(w, http.StatusBadRequest, "Invalid Action")
		return
	}
}
