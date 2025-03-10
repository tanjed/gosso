package auth

import (
	"net/http"

	"github.com/tanjed/go-sso/internal/customerror"
	"github.com/tanjed/go-sso/internal/handler/customtype"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/pkg/responsemanager"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var userRegisterRequest customtype.UserRegisterRequest
	
	if err := userRegisterRequest.Validated(r.Body); err != nil {
		responsemanager.ResponseUnprocessableEntity(&w, customtype.I{
			"message" : err.Error(),
			"bag" : err.(*customerror.ValidationError).ErrBag,
		})
        return
	}

	if _, err := model.GetClientById(userRegisterRequest.ClientId); err != nil {
		responsemanager.ResponseUnAuthorized(&w, customtype.M{"message" : "invalid client"})
		return
	}

	exists, err := model.GetUserByMobileNumber(userRegisterRequest.MobileNumber);
	if err != nil && err != mongo.ErrNoDocuments{
		responsemanager.ResponseServerError(&w, customtype.M{"message" : "something went wrong"})
		return
	}

	if exists != nil {
		responsemanager.ResponseUnprocessableEntity(&w, customtype.M{"message" : "user already exists"})
		return
	}

	user := model.NewUser(userRegisterRequest)
	if !user.Insert() {
		responsemanager.ResponseServerError(&w, customtype.M{"message" : "unable to create user"})
		return
	}

	responsemanager.ResponseOK(&w, customtype.I{
		"user" : customtype.I{
			"user_id" : user.UserId,
			"first_name" : user.FirstName,
			"last_name" : user.LastName,
			"mobile_number" : user.MobileNumber,
			"email" : user.Email,
		},
	})
}