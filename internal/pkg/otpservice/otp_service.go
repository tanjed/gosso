package otpservice

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/tanjed/go-sso/internal/customerror"
	"github.com/tanjed/go-sso/internal/model"
)


func SendResetOtp(userId string, otpType string) error{
	otp := model.NewOtp(userId, GenerateOTP(), otpType)
	if !otp.Insert() {
		return errors.New("unable to store otp")
	}

	//Process sending otp with api call
	return nil
}

func GenerateOTP() string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	otp := r.Intn(900000) + 100000
	return fmt.Sprintf("%d", otp)
}

func GeneratePasswordResetToken() string {
	length := 35
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}


func IsReadyToSendOtp(user *model.User, requestType string, isResend bool) (error){
	if user == nil {
		return &customerror.UserNotFoundError{
			ErrMessage: "user not found",
			ErrCode: 404,
		}
	}

	sentOtpCount, err := model.GetUserSentOtpCount(user.UserId, requestType)

	if err != nil {
		return &customerror.DBReadError{
			ErrMessage: "unable to read user otp count",
			ErrCode: http.StatusInternalServerError, 
		}
	}

	if sentOtpCount >= model.MAX_RETRY_COUNT {
		return &customerror.OtpLimitReachedError{
			ErrMessage: "max retry limit reached",
			ErrCode: http.StatusTooManyRequests,
		}
	}

	validSentOtp, _ := model.GetUserValidOtp(user.UserId, model.OTP_TYPE_PASSWORD_RESET)

	if !isResend && validSentOtp != nil {	
		return &customerror.OtpLimitReachedError{
			ErrMessage: "otp already sent",
			ErrCode: http.StatusTooManyRequests,
		}
	}

	return nil
}


func ValidateOtp(userId string, requestOtp string, otpType string) (*model.PasswordReset, error) {
	otp, err := model.GetUserValidOtp(userId, otpType)

	if err != nil {
		return nil, &customerror.OtpNotFoundError{
			ErrMessage: "otp not found",
			ErrCode: http.StatusUnprocessableEntity,
		}
	}

	if otp.Otp != requestOtp {
		return nil, &customerror.OtpMismatchError{
			ErrMessage: "otp mismatched",
			ErrCode: http.StatusUnprocessableEntity,
		}
	}

	passwordReset := model.NewPasswordReset(userId, GeneratePasswordResetToken())
	
	if !otp.MarkAsValidated() || !passwordReset.Insert() {
		return nil, &customerror.OtpMismatchError{
			ErrMessage: "something went wrong",
			ErrCode: http.StatusInternalServerError,
		}
	}

	return passwordReset, nil
}