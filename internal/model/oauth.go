package model

import (
	"log/slog"
	"time"

	"github.com/tanjed/go-sso/internal/db"
)
const TOKEN_TYPE_USER_ACCESS_TOKEN = "USER_ACCESS_TOKEN"
const TOKEN_TYPE_USER_REFRESH_TOKEN = "USER_REFRESH_TOKEN"
const TOKEN_TYPE_CLIENT_ACCESS_TOKEN = "CLIENT_ACCESS_TOKEN"
const TOKEN_TYPE_CLIENT_REFRESH_TOKEN = "CLIENT_REFRESH_TOKEN"

type TokenableInterface interface {
	Insert() bool
	InvokeToken() bool
	GetExpiry() time.Time
	GetClientId() string
	GetScopes() []string
}

type Token struct {
	TokenId string
	ClientId string
	Scopes []string
	Revoked int
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserAccessToken struct {
	UserId string
	Token
}

type UserRefreshToken struct {
	UserId string
	Token
}

type ClientAccessToken struct {
	Token
}

type ClientRefreshToken struct {
	Token
}



func (t *UserAccessToken) Insert() bool {
	db := db.InitDB()
	
	err := db.Conn.Query("INSERT INTO user_access_tokens (token_id, client_id, user_id, scopes, revoked, expired_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", 
	t.TokenId, t.ClientId, t.UserId, t.Scopes, t.Revoked, t.ExpiredAt, t.CreatedAt, t.UpdatedAt).Exec()

	if err != nil {
		slog.Error("Unable to record token", "error", err)
		return false
	}

	return true
}

func (t *UserAccessToken) GetExpiry() time.Time {
	return t.ExpiredAt
}

func (t *UserAccessToken) GetClientId() string {
	return t.ClientId
}

func (t *UserAccessToken) GetScopes() []string {
	return t.Scopes
}

func (t *UserAccessToken) InvokeToken() bool{
	db := db.InitDB()
		
	err := db.Conn.Query("UPDATE user_access_tokens SET revoked = 1, updated_at = ? WHERE token_id = ? AND user_id = ?", 
	time.Now(),t.TokenId, t.UserId).Exec()

	if err != nil {
		slog.Error("Unable to invoke token", "error", err)
		return false
	}

	return true
}


func (t *UserRefreshToken) Insert() bool {
	db := db.InitDB()
	
	err := db.Conn.Query("INSERT INTO user_refresh_tokens (token_id, client_id, user_id, scopes, revoked, expired_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", 
	t.TokenId, t.ClientId, t.UserId, t.Scopes, t.Revoked, t.ExpiredAt, t.CreatedAt, t.UpdatedAt).Exec()

	if err != nil {
		slog.Error("Unable to record token", "error", err)
		return false
	}

	return true
}

func (t *UserRefreshToken) GetExpiry() time.Time {
	return t.ExpiredAt
}

func (t *UserRefreshToken) GetProperties() *UserRefreshToken {
	return t
}

func (t *UserRefreshToken) GetClientId() string {
	return t.ClientId
}

func (t *UserRefreshToken) GetScopes() []string {
	return t.Scopes
}


func (t *UserRefreshToken) InvokeToken() bool{
	db := db.InitDB()
		
	err := db.Conn.Query("UPDATE user_refresh_tokens SET revoked = 1, updated_at = ? WHERE token_id = ? AND user_id = ?", 
	time.Now(), t.TokenId, t.UserId).Exec()

	if err != nil {
		slog.Error("Unable to invoke token", "error", err)
		return false
	}

	return true
}


func (t *ClientAccessToken) Insert() bool {
	db := db.InitDB()
	
	err := db.Conn.Query("INSERT INTO client_access_tokens (token_id, client_id, scopes, revoked, expired_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)", 
	t.TokenId, t.ClientId, t.Scopes, t.Revoked, t.ExpiredAt, t.CreatedAt, t.UpdatedAt).Exec()

	if err != nil {
		slog.Error("Unable to record token", "error", err)
		return false
	}

	return true
}

func (t *ClientAccessToken) GetExpiry() time.Time {
	return t.ExpiredAt
}

func (t *ClientAccessToken) GetProperties() *ClientAccessToken {
	return t
}


func (t *ClientAccessToken) GetClientId() string {
	return t.ClientId
}

func (t *ClientAccessToken) GetScopes() []string {
	return t.Scopes
}

func (t *ClientAccessToken) InvokeToken() bool{
	db := db.InitDB()
		
	err := db.Conn.Query("UPDATE client_access_tokens SET revoked = 1, updated_at = ? WHERE token_id = ? AND client_id = ?", 
	time.Now(), t.TokenId, t.ClientId).Exec()

	if err != nil {
		slog.Error("Unable to invoke token", "error", err)
		return false
	}

	return true
}

func (t *ClientRefreshToken) Insert() bool {
	db := db.InitDB()
	
	err := db.Conn.Query("INSERT INTO client_refresh_tokens (token_id, client_id, scopes, revoked, expired_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)", 
	t.TokenId, t.ClientId, t.Scopes, t.Revoked, t.ExpiredAt, t.CreatedAt, t.UpdatedAt).Exec()

	if err != nil {
		slog.Error("Unable to record token", "error", err)
		return false
	}

	return true
}

func (t *ClientRefreshToken) GetExpiry() time.Time {
	return t.ExpiredAt
}

func (t *ClientRefreshToken) GetProperties() *ClientRefreshToken {
	return t
}

func (t *ClientRefreshToken) GetClientId() string {
	return t.ClientId
}

func (t *ClientRefreshToken) GetScopes() []string {
	return t.Scopes
}

func (t *ClientRefreshToken) InvokeToken() bool{
	db := db.InitDB()
		
	err := db.Conn.Query("UPDATE client_refresh_tokens SET revoked = 1, updated_at = ? WHERE token_id = ? AND client_id = ?", 
	time.Now(), t.TokenId, t.ClientId).Exec()

	if err != nil {
		slog.Error("Unable to invoke token", "error", err)
		return false
	}

	return true
}

func GetOAuthTokenById(tokenId string, model TokenableInterface) TokenableInterface {
	db := db.InitDB()

	if clientAccessToken, ok := model.(*ClientAccessToken); ok {
		if err := db.Conn.Query("SELECT token_id, client_id, revoked, expired_at, created_at, updated_at FROM client_access_tokens WHERE token_id = ?", tokenId).
		Scan(
			&clientAccessToken.TokenId,
			&clientAccessToken.ClientId,
			&clientAccessToken.Revoked,
			&clientAccessToken.ExpiredAt,
			&clientAccessToken.CreatedAt,
			&clientAccessToken.UpdatedAt,
		); err != nil {
			slog.Error("Unable to fetch token", "error", err)
			return nil
		}
		return clientAccessToken
		
	} else if clientRefreshToken, ok := model.(*ClientRefreshToken); ok {
		if err := db.Conn.Query("SELECT token_id, client_id, revoked, expired_at, created_at, updated_at FROM client_refresh_tokens WHERE token_id = ?", tokenId).
		Scan(
			&clientRefreshToken.TokenId,
			&clientRefreshToken.ClientId,
			&clientRefreshToken.Revoked,
			&clientRefreshToken.ExpiredAt,
			&clientRefreshToken.CreatedAt,
			&clientRefreshToken.UpdatedAt,
		); err != nil {
			slog.Error("Unable to fetch token", "error", err)
			return nil
		}
		return clientRefreshToken
	} else if userAccessToken, ok := model.(*UserAccessToken); ok {
		if err := db.Conn.Query("SELECT token_id, client_id, user_id, revoked, expired_at, created_at, updated_at FROM user_access_tokens WHERE token_id = ?", tokenId).
		Scan(
			&userAccessToken.TokenId,
			&userAccessToken.ClientId,
			&userAccessToken.UserId,
			&userAccessToken.Revoked,
			&userAccessToken.ExpiredAt,
			&userAccessToken.CreatedAt,
			&userAccessToken.UpdatedAt,
		); err != nil {
			slog.Error("Unable to fetch token", "error", err)
			return nil
		}
		return userAccessToken
		
	} else if userRefreshToken, ok := model.(*UserRefreshToken); ok {
		if err := db.Conn.Query("SELECT token_id, client_id, user_id, revoked, expired_at, created_at, updated_at FROM user_refresh_tokens WHERE token_id = ?", tokenId).
		Scan(
			&userRefreshToken.TokenId,
			&userRefreshToken.ClientId,
			&userRefreshToken.UserId,
			&userRefreshToken.Revoked,
			&userRefreshToken.ExpiredAt,
			&userRefreshToken.CreatedAt,
			&userRefreshToken.UpdatedAt,
		); err != nil {
			slog.Error("Unable to fetch token", "error", err)
			return nil
		}
		return userRefreshToken
	} else {
		return nil
	}
}




