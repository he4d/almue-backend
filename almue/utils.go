package almue

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
