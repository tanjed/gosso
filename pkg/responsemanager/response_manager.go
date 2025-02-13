package responsemanager

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func ResponseOK(w *http.ResponseWriter, data interface{}) {
	responseWriter := *w
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(map[string]interface{}{
		"data" : data,
	})
}

func ResponseUnprocessableEntity(w *http.ResponseWriter, message string) {
	responseWriter := *w
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(responseWriter).Encode(map[string]interface{}{
		"errors" : message,
	})
}


func ResponseUnAuthorized(w *http.ResponseWriter, message string) {
	responseWriter := *w
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(responseWriter).Encode(map[string]interface{}{
		"errors" : message,
	})
}

func ResponseValidationError(w *http.ResponseWriter, errors validator.ValidationErrors) {
	responseWriter := *w
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusUnprocessableEntity)
	validationErrors := make(map[string][]string)
	for _, fieldError := range errors {
		validationErrors[fieldError.Field()] = append(
			validationErrors[fieldError.Field()], 
			fmt.Sprintf("'%s' failed validation: %s", fieldError.Field(), fieldError.Tag())) 
	}

	json.NewEncoder(responseWriter).Encode(map[string]interface{}{
		"errors" : validationErrors,
	})
}