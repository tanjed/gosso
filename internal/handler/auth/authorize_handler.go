package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gocql/gocql"
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
	
	client, err := model.AuthenticateClient(authenticationRequest.ClientName, authenticationRequest.ClientSecret);
	if err != nil {
		responsemanager.ResponseUnAuthorized(&w, "Invalid client credentials")
		return
	}

	if _, err := model.AutheticateUser(authenticationRequest.MobileNumber, authenticationRequest.Password); err != nil  {
		responsemanager.ResponseUnAuthorized(&w, "Invalid user credentials")
		return
	}

	clientAuthorizationCode, err := model.ClientHasValidSession(client.ClientId, authenticationRequest.RedirectUri)
	if err != nil {
		if err != gocql.ErrNotFound{
			responsemanager.ResponseServerError(&w, "Something went wrong")
			return
		} 
		
		clientAuthorizationCode, err = model.NewClientAuthorizationCode(client.ClientId, authenticationRequest.RedirectUri)
		if err != nil{
			responsemanager.ResponseServerError(&w, "Something went wrong")
			return
		} 
	}

	responsemanager.ResponseOK(&w, map[string]interface{}{
		"code" : clientAuthorizationCode.ClientCode,
		"redirect_uri" : clientAuthorizationCode.RedirectUri,
		"expired_at" : clientAuthorizationCode.ExpiredAt,
	})
}