package hashutilities

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHashFromString(s string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
    if err != nil {
        log.Println(err)
        return ""
    }

	return string(hashedPassword)
}

func CompareHashWithString(h string, s string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(s))
	if err != nil {
        log.Println(err)
        return false
    }

	return true
}