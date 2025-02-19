package model

import (
	"log/slog"
	"time"

	"github.com/tanjed/go-sso/internal/db"
)

type OauthToken struct {
	TokenId string
	ClientId string
	UserId string
	Scopes []string
	Revoked int
	Type string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}


func (t *OauthToken) Insert() bool {
	db := db.InitDB()
		
	err := db.Conn.Query("INSERT INTO oauth_tokens (token_id, client_id, user_id, scopes, revoked, type, expired_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", 
	t.TokenId, t.ClientId, t.UserId, t.Scopes, t.Revoked, t.Type, t.ExpiredAt, t.CreatedAt, t.UpdatedAt).Exec()

	if err != nil {
		slog.Error("Unable to record token", "error", err)
		return false
	}

	return true
}

func NewOauthToken(tokenId, clientId, userId string, scopes []string, revoked int, tokenType string, expiredAt, created_at, updated_at time.Time) *OauthToken{
	return &OauthToken{
		TokenId: tokenId,
		ClientId: clientId,
		UserId: userId,
		Scopes: scopes,
		Revoked: revoked,
		Type: tokenType,
		ExpiredAt: expiredAt,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}
}

