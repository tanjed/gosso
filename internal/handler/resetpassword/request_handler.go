package resetpassword

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/internal/pkg/otpservice"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)

type ResetRequest struct {
	IsResend bool `json:"is_resend,omitempty"`
}


func RequestHandler(w http.ResponseWriter, r *http.Request) {
	var resetRequest ResetRequest
	
	if r.ContentLength != 0 {
		if err:= json.NewDecoder(r.Body).Decode(&resetRequest); err != nil {
			slog.Error("Unable to process the request", "error", err)
			responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
			return
		}
	
		validate := validator.New()
	
		if err := validate.Struct(resetRequest); err != nil {
			if errors, ok := err.(validator.ValidationErrors); ok {
				responsemanager.ResponseValidationError(&w, errors)
			} else {
				slog.Error("Unable to process the request", "error", err)
				responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
			}
			return
		}	
	}	

	user := r.Context().Value(model.AUTH_USER_CONTEXT_KEY).(*model.User)

	if err := otpservice.IsReadyToSendOtp(user, model.OTP_TYPE_PASSWORD_RESET, resetRequest.IsResend); err != nil {
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