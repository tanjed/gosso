package forgotpassword

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/internal/pkg/otpservice"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)

type ValidateRequest struct {
	MobileNumber string `json:"mobile_number"`
	Otp string `json:"otp"`
}

func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	var validateRequest ValidateRequest
	if err:= json.NewDecoder(r.Body).Decode(&validateRequest); err != nil {
		slog.Error("Unable to process the request", "error", err)
		responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
        return
	}

	validate := validator.New()
	
	if err := validate.Struct(validateRequest); err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			responsemanager.ResponseValidationError(&w, errors)
		} else {
			slog.Error("Unable to process the request", "error", err)
			responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
		}
		return
    }
	

	user := model.GetUserByMobileNumber(validateRequest.MobileNumber)

	if user == nil {
		responsemanager.ResponseUnAuthorized(&w, "otp can not be verified")
		return
	}

	passwordReset, err := otpservice.ValidateOtp(user.UserId, validateRequest.Otp, model.OTP_TYPE_PASSWORD_RESET)

	if err != nil {
		responsemanager.ResponseWithCode(&w, err)
		return
	}
	
	responsemanager.ResponseOK(&w, map[string]interface{}{
		"token" : passwordReset.Token,
		"expires_at" : passwordReset.ExpiredAt.Unix(),
	})
}
