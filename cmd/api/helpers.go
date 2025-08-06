package main

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, data any) {
	payload, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "could not marshal json response", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(payload)
}
