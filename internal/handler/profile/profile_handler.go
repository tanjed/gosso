package profile

import (
	"net/http"

	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	authUser := r.Context().Value(model.AUTH_USER_CONTEXT_KEY);
	var responsePayload map[string]interface{}
	if client, ok := authUser.(*model.Client); ok {
		responsePayload = map[string]interface{}{
			"client_id" : client.ClientId,
			"client_name" : client.ClientName,
		}
	} 

	if user, ok := authUser.(*model.User); ok {
		responsePayload = map[string]interface{}{
			"first_name" : user.FirstName,
			"last_name" : user.LastName,
			"mobile_number" : user.MobileNumber,
		}
	} 
	responsemanager.ResponseOK(&w, responsePayload)
}