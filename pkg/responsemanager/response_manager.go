package responsemanager

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/tanjed/go-sso/internal/customerror"
)

func ResponseOK(w *http.ResponseWriter, data interface{}) {
	responseWriter := *w
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(map[string]interface{}{
		"success" : true,
		"data" : data,
	})
}

func ResponseUnprocessableEntity(w *http.ResponseWriter, d interface{}) {
	responseWriter := *w
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(responseWriter).Encode(map[string]interface{}{
		"success" : false,
		"errors" : d,
	})
}


func ResponseServerError(w *http.ResponseWriter, message interface{}) {
	responseWriter := *w
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(responseWriter).Encode(map[string]interface{}{
		"success" : false,
		"errors" : message,
	})
}


func ResponseUnAuthorized(w *http.ResponseWriter, message interface{}) {
	responseWriter := *w
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(responseWriter).Encode(map[string]interface{}{
		"success" : false,
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
		"success" : false,
		"errors" : validationErrors,
	})
}


func ResponseWithCode(w *http.ResponseWriter, err error) {
	if errableElement, ok := err.(customerror.ErrorableInterface); ok {
		switch errableElement.Code() {
		case http.StatusUnprocessableEntity :
			ResponseUnprocessableEntity(w, errableElement.Message())
		default :
			ResponseServerError(w, errableElement.Message())
		}
	}else {
		ResponseServerError(w, err.Error())
	}
}