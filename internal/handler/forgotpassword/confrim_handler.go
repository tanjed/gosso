package forgotpassword

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)

type ConfirmRequest struct {
	Token string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required,eqfield=NewPassword"`
}

func ConfirmHandler(w http.ResponseWriter, r *http.Request) {
	var confirmRequest ConfirmRequest
	if err:= json.NewDecoder(r.Body).Decode(&confirmRequest); err != nil {
		slog.Error("Unable to process the request", "error", err)
		responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
        return
	}

	validate := validator.New()

	if err := validate.Struct(confirmRequest); err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			responsemanager.ResponseValidationError(&w, errors)
		} else {
			slog.Error("Unable to process the request", "error", err)
			responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
		}
		return
    }
	
	reset, err := model.GetResetTokenById(confirmRequest.Token)

	if err != nil {
		responsemanager.ResponseUnprocessableEntity(&w, "invalid reset token")
		return
	}

	if reset == nil || reset.Token != confirmRequest.Token {
		responsemanager.ResponseUnprocessableEntity(&w, "invalid token")
		return
	}	
	
	u, _ := model.GetUserById(reset.UserId)
	
	if !reset.MarkAsValidated() && !u.UpdatePassword(confirmRequest.NewPassword) {
		responsemanager.ResponseUnprocessableEntity(&w, "password not updated")
		return
	}

	responsemanager.ResponseOK(&w, "password updated")
}
