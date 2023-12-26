package api

import (
	"database/sql"
	"net/http"
	"strings"

	db "github.com/broemp/red_card/db/sqlc"
	"github.com/broemp/red_card/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username  string `json:"username" binding:"required,alphanum"`
	Password  string `json:"password" binding:"required,min=8,max=256"`
	FirstName string `json:"first_name" binding:"max=64" `
	LastName  string `json:"last_name" binding:"max=64"`
}

type userResponse struct {
	ID        int64          `json:"id"`
	Username  string         `json:"username"`
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	CreatedAt sql.NullTime   `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	firstName := util.StringToSQLString(req.FirstName)
	lastName := util.StringToSQLString(req.LastName)

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       strings.ToLower(req.Username),
		HashedPassword: hashedPassword,
		FirstName:      firstName,
		LastName:       lastName,
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
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=256"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.store.GetUser(ctx, strings.ToLower(req.Username))
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

	accessToken, err := server.tokenMaker.CreateToken(user, server.config.JWT_DURATION)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
