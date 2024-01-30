package tools

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Message string      `json:"error"`
	Payload interface{} `json:"payload,omitempty"`
}

func WriteJSON(w http.ResponseWriter, data interface{}, headers ...http.Header) error {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJSON(w http.ResponseWriter, err error, status int) error {
	w.WriteHeader(status)

	errorPayload := JSONResponse{
		Message: err.Error(),
	}

	return WriteJSON(w, errorPayload)
}
