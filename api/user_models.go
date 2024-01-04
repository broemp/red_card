package api

import (
	"time"

	db "github.com/broemp/red_card/db/sqlc"
	"github.com/google/uuid"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=256"`
	Name     string `json:"name" binding:"required,max=64" `
}

type userResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	ID        int64     `json:"id"`
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=256"`
}

type loginUserResponse struct {
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	SessionID             uuid.UUID    `json:"session_id"`
	RefressTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	AccessToken           string       `json:"access_token"`
	RefreshToken          string       `json:"refresh_token"`
	User                  userResponse `json:"user"`
}

type listUserFilterRequest struct {
	Filter   string `form:"filter" `
	PageID   int32  `form:"page_id" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=5,max=100"`
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type getUserCardsRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type getUserCardsResponse struct {
	Cards []db.ListCardsByUserIDRow         `json:"cards"`
	Count []db.GetCardColorCountByUserIDRow `json:"count"`
}
