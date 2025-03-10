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

type ForgotRequest struct {
	MobileNumber string `mobile_number:"mobile_number"`
	IsResend bool `json:"is_resend,omitempty"`
}

func RequestHandler(w http.ResponseWriter, r *http.Request){
	var forgotPasswordRequest ForgotRequest
	if err:= json.NewDecoder(r.Body).Decode(&forgotPasswordRequest); err != nil {
		slog.Error("Unable to process the request", "error", err)
		responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
		return
	}

	validate := validator.New()

	if err := validate.Struct(forgotPasswordRequest); err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			responsemanager.ResponseValidationError(&w, errors)
		} else {
			slog.Error("Unable to process the request", "error", err)
			responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
		}
		return
	}
	
	user, err := model.GetUserByMobileNumber(forgotPasswordRequest.MobileNumber)
	if err != nil {
		responsemanager.ResponseServerError(&w, "something went wrong")
		return
	}

	if err := otpservice.IsReadyToSendOtp(user, model.OTP_TYPE_PASSWORD_RESET, forgotPasswordRequest.IsResend); err != nil {
		responsemanager.ResponseWithCode(&w, err)
		return
	}

	if err := otpservice.SendResetOtp(user.UserId, model.OTP_TYPE_PASSWORD_RESET); err != nil {
		responsemanager.ResponseServerError(&w, "something went wrong")
		return
	}
	
	responsemanager.ResponseOK(&w, map[string]interface{}{
		"message" : "otp sent",
	})

}