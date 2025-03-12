package profile

import (
	"net/http"

	"github.com/tanjed/go-sso/internal/handler/customtype"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(model.AUTH_USER_CONTEXT_KEY);
	
	if _, ok:= user.(model.Client); ok {
		responsemanager.ResponseUnprocessableEntity(&w, customtype.M{"message" : "client does not have this permission"})
		return
	}

	responsemanager.ResponseOK(&w, customtype.M{
		"first_name" : user.(*model.User).FirstName,
		"last_name" : user.(*model.User).LastName,
		"mobile_number" : user.(*model.User).MobileNumber,
		"email" : user.(*model.User).Email,
		"address" : *user.(*model.User).Address,
		"nid" : *user.(*model.User).Nid,
		"passport" : *user.(*model.User).Passport,
	})
}