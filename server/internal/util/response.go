package util

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/madeindra/mock-interview/server/internal/model"
)

func SendResponse(w http.ResponseWriter, data any, message string, status int) {
	resp, err := json.Marshal(model.Response{
		Message: message,
		Data:    data,
	})
	if err != nil {
		log.Printf("failed to marshal response: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "an error occured while processing the request"}`))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}
