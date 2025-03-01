package otpservice

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

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