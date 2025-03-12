package customtype

import (
	"io"

	"github.com/tanjed/go-sso/internal/handler/request"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type M map[string]string
type I map[string]interface{}

type TokenRequest struct {
	// ClientId bson.ObjectID `json:"client_id" validate:"required"`
	GrantType string 	`json:"grant_type" validate:"required"`
}
type UserRegisterRequest struct {
	ClientId bson.ObjectID `json:"client_id" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName string `json:"last_name" validate:"required"`
	MobileNumber string `json:"mobile_number" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Address string `json:"address,omitempty"`
	Nid string `json:"nid,omitempty"`
	Passport string `json:"passport,omitempty"`
}

func(r *UserRegisterRequest) Validated(d io.Reader) error {
	return request.GetValidated(d, r)
}


type PasswordGrantRequest struct {
	ClientId bson.ObjectID 	`json:"client_id" validate:"required"`
	Scope []string 		`json:"scope" validate:"required"`
	MobileNumber string `json:"mobile_number" validate:"required"`
	Password string 	`json:"password" validate:"required"`
}

func(r *PasswordGrantRequest) Validated(d io.Reader) error {
	return request.GetValidated(d, r)
}


type ClientCredentialsGrantRequest struct {
	TokenRequest
	ClientId bson.ObjectID 	`json:"client_id" validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
	Scope []string 		`json:"scope" validate:"required"`
}


func(r *ClientCredentialsGrantRequest) Validated(d io.Reader) error {
	return request.GetValidated(d, r)
}


type RefreshTokenRequest struct {
	TokenRequest
	AccessToken string `json:"access_token" validate:"required"`
}

func(r *RefreshTokenRequest) Validated(d io.Reader) error {
	return request.GetValidated(d, r)
}

