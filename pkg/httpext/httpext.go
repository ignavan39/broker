package httpext

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, body any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func AbortJSON(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(CommonError{
		Error: err,
	})
}
