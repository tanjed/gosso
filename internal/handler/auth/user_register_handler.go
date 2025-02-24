package auth

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/internal/validation"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)

type UserRegisterRequest struct {
	ClientName string `json:"client_name" validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName string `json:"last_name" validate:"required"`
	MobileNumber string `json:"mobile_number" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var userRegisterRequest UserRegisterRequest
	
	if err := json.NewDecoder(r.Body).Decode(&userRegisterRequest); err != nil {
		slog.Error("Unable to process the request", "error", err)
		responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
        return
	}

	validate := validator.New()
	validate.RegisterValidation("mobileno", validation.ValidateMobile)
	if err := validate.Struct(userRegisterRequest); err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			responsemanager.ResponseValidationError(&w, errors)
		} else {
			slog.Error("Unable to process the request", "error", err)
			responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
		}
		return
    }

	if _, err := model.AuthenticateClient(userRegisterRequest.ClientName, userRegisterRequest.ClientSecret); err != nil {
		responsemanager.ResponseUnAuthorized(&w, "Invalid client credentials")
		return
	}

	if existingUser := model.GetUserByMobileNumber(userRegisterRequest.MobileNumber); existingUser != nil {
		log.Println(existingUser)
		responsemanager.ResponseUnprocessableEntity(&w,"User already exists")
		return
	}

	user := model.NewUser(userRegisterRequest.FirstName, userRegisterRequest.LastName, userRegisterRequest.MobileNumber, userRegisterRequest.Password)
	createdUser := user.Insert()

	if createdUser == nil {
		responsemanager.ResponseServerError(&w, "Unable to create user")
		return
	}

	responsemanager.ResponseOK(&w, map[string]interface{}{
		"user" : map[string]interface{}{
			"user_id" : createdUser.UserId,
			"first_name" : createdUser.FirstName,
			"last_name" : createdUser.LastName,
			"mobile_number" : createdUser.MobileNumber,
		},
	})
}