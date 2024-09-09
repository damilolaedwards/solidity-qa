package handlers

import (
	"log"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	log.Println(message)
	w.WriteHeader(statusCode)
}
