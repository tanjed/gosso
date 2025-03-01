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
const MAX_RETRY_COUNT = 3

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

	if user == nil {
		responsemanager.ResponseOK(&w, map[string]interface{}{
			"message" : "otp sent",
		})
		return
	}

	sentOtpCount, err := model.GetUserSentOtpCount(user.UserId, model.OTP_TYPE_PASSWORD_RESET)

	if err != nil {
		responsemanager.ResponseServerError(&w, "something went wrong")
		return
	}

	if sentOtpCount >= MAX_RETRY_COUNT {
		responsemanager.ResponseUnprocessableEntity(&w, "max retry limit reached")
		return
	}

	validSentOtp, _ := model.GetUserValidOtp(user.UserId, model.OTP_TYPE_PASSWORD_RESET)

	if !resetRequest.IsResend && validSentOtp != nil {	
		responsemanager.ResponseUnprocessableEntity(&w, "otp already sent")
		return
	}

	err = otpservice.SendResetOtp(user.UserId, model.OTP_TYPE_PASSWORD_RESET)

	if err != nil {
		responsemanager.ResponseServerError(&w, "something went wrong")
		return
	}
	
	responsemanager.ResponseOK(&w, map[string]interface{}{
		"message" : "otp sent",
	})

}