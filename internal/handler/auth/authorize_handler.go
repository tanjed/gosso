package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)

type AuthenticationRequest struct {
	ClientName string 	`json:"client_name" validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
	MobileNumber string `json:"mobile_number" validate:"required"`
	Password string 	`json:"password" validate:"required"`
	RedirectUri string 	`json:"redirect_uri" validate:"required"`
	State string 		`json:"state" validate:"required"`
	Scope string 		`json:"scope" validate:"required"`
	ResponseType string `json:"response_type" validate:"required"`
}

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	var authenticationRequest AuthenticationRequest
	err := json.NewDecoder(r.Body).Decode(&authenticationRequest); 

	if err != nil {
		log.Println(err)
		responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
        return
	}

	validate := validator.New()
	err = validate.Struct(authenticationRequest)

	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			responsemanager.ResponseValidationError(&w, errors)
		} else {
			responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
		}
		return
    }

	if !model.AuthenticateClient(authenticationRequest.ClientName, authenticationRequest.ClientSecret) {
		responsemanager.ResponseUnAuthorized(&w, "Invalid client credentials")
		return
	}

	responsemanager.ResponseOK(&w, "User created successfully!")
}