package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/tanjed/go-sso/internal/handler/customtype"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)

const GRANT_TYPE_PASSWORD = "password"
const GRANT_TYPE_REFRESH_TOKEN = "refresh_token"
const GRANT_TYPE_CLIENT_CREDENTIALS = "client_credentials"


func TokenHandler(w http.ResponseWriter, r *http.Request) {
	
	var tokenRequest customtype.TokenRequest 
	body, err := io.ReadAll(r.Body)

	if err != nil {
		slog.Error("unable to process the request", "error", err)
		responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
        return
	}

	err = json.Unmarshal(body, &tokenRequest)
	
	if err != nil {
		slog.Error("unable to process the request", "error", err)
		responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
        return
	}

	r.Body = io.NopCloser(bytes.NewReader(body)) //Reseting the request body to its original form

	switch tokenRequest.GrantType {
		case GRANT_TYPE_PASSWORD :
			passwordGrantHandler(w, r)
		case GRANT_TYPE_CLIENT_CREDENTIALS :
			clientCredentialsGrantHandler(w, r)
		case GRANT_TYPE_REFRESH_TOKEN : 
			refreshTokenGrantHandler(w, r)
		default : 
			responsemanager.ResponseUnprocessableEntity(&w, "Invalid grant type")
	}

}