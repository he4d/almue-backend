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

func (a *Almue) getAllShuttersOfFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i, err := strconv.ParseInt(vars["floorID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid floor ID")
		return
	}
	shutters, err := a.store.GetShutterListOfFloor(i)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, shutters)
}

func (a *Almue) getShutter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	floorID, err := strconv.ParseInt(vars["floorID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid floor ID")
		return
	}
	shutterID, err := strconv.ParseInt(vars["shutterID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid shutter ID")
		return
	}
	shutter, err := a.store.GetShutterByFloor(shutterID, floorID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Shutter not found")
			return
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, shutter)
}

func (a *Almue) createShutter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i, err := strconv.ParseInt(vars["floorID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid floor ID")
		return
	}
	s := new(model.Shutter)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	s.FloorID = i

	_, err = govalidator.ValidateStruct(s)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	newID, err := a.store.CreateShutter(s)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	shutter, err := a.store.GetShutterByFloor(newID, i)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	shutterState, err := a.deviceController.RegisterShutters(shutter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = a.registerShutterStateSynchronization(newID, shutterState[newID]); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, shutter)
}

func (a *Almue) updateShutter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	floorID, err := strconv.ParseInt(vars["floorID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid floor ID")
		return
	}

	shutterID, err := strconv.ParseInt(vars["shutterID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid shutter ID")
		return
	}

	s := new(model.Shutter)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()

	_, err = govalidator.ValidateStruct(s)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.store.UpdateShutter(s); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	shutter, err := a.store.GetShutterByFloor(shutterID, floorID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//TODO: update a.devices

	respondWithJSON(w, http.StatusOK, shutter)
}

func (a *Almue) deleteShutter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shutterID, err := strconv.ParseInt(vars["shutterID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid shutter ID")
		return
	}

	if err := a.store.DeleteShutter(shutterID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := a.deviceController.UnregisterShutter(shutterID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (a *Almue) controlShutter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	floorID, err := strconv.ParseInt(vars["floorID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid floor ID")
		return
	}
	shutterID, err := strconv.ParseInt(vars["shutterID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid lighting ID")
		return
	}

	_, err = a.store.GetShutterByFloor(shutterID, floorID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Shutter not found")
		return
	}

	action := vars["action"]
	switch action {
	case "open":
		if err := a.deviceController.OpenShutter(shutterID); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		break
	case "close":
		if err := a.deviceController.CloseShutter(shutterID); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		break
	case "stop":
		if err := a.deviceController.StopShutter(shutterID); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		break
	default:
		respondWithError(w, http.StatusBadRequest, "Invalid Action")
		return
	}
}
