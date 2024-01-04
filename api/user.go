package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	db "github.com/broemp/red_card/db/sqlc"
	"github.com/broemp/red_card/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=256"`
	Name     string `json:"name" binding:"required,max=64" `
}

type userResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       strings.ToLower(req.Username),
		HashedPassword: hashedPassword,
		Name:           req.Name,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusCreated, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=256"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefressTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.store.GetUserAuth(ctx, strings.ToLower(req.Username))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accesPayload, err := server.tokenMaker.CreateToken(user, server.config.JWT_DURATION)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user, server.config.JWT_REFRESH_DURATION)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payloadID, err := uuid.Parse(refreshPayload.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           payloadID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accesPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefressTokenExpiresAt: refreshPayload.ExpiresAt.Time,
		User:                  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}

type listUserFilterRequest struct {
	Filter   string `form:"filter" `
	PageID   int32  `form:"page_id" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=5,max=100"`
}

func (server *Server) listUserFilter(ctx *gin.Context) {
	var req listUserFilterRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUserFilterParams{
		Filter: req.Filter,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUserFilter(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		fmt.Println(req)

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserCardsRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type getUserCardsResponse struct {
	Cards []db.ListCardsByUserIDRow         `json:"cards"`
	Count []db.GetCardColorCountByUserIDRow `json:"count"`
}

func (server *Server) getUserCards(ctx *gin.Context) {
	var req getUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		fmt.Println(req)

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cards, err := server.store.ListCardsByUserID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	count, err := server.store.GetCardColorCountByUserID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := getUserCardsResponse{
		Cards: cards,
		Count: count,
	}

	ctx.JSON(http.StatusOK, rsp)
}
