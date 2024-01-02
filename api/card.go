package api

import (
	"database/sql"
	"net/http"
	"strconv"

	db "github.com/broemp/red_card/db/sqlc"
	"github.com/broemp/red_card/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createCardRequest struct {
	Color       string `json:"color" binding:"required,oneof=red yellow blue"`
	Accused     int64  `json:"accused" binding:"required"`
	Description string `json:"description"`
	Event       int64  `json:"event"`
}

func (s *Server) createCard(ctx *gin.Context) {
	var req createCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Claims)

	var color db.Color
	switch req.Color {
	case "red":
		color = db.ColorRed
	case "yellow":
		color = db.ColorYellow
	case "blue":
		color = db.ColorBlue
	}

	// accusedID, err := s.store.GetUserID(ctx, req.Accused)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }

	// Parse Subject string to Int64
	authorID, err := strconv.ParseInt(authPayload.Subject, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	description := sql.NullString{Valid: false}
	if req.Description != "" {
		description = sql.NullString{
			Valid:  true,
			String: req.Description,
		}
	}

	event := sql.NullInt64{Valid: false}
	if req.Event != 0 {
		event = sql.NullInt64{
			Valid: true,
			Int64: req.Event,
		}
	}

	arg := db.CreateCardParams{
		Author:      authorID,
		Accused:     req.Accused,
		Color:       color,
		Event:       event,
		Description: description,
	}

	card, err := s.store.CreateCard(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			case "unique_violation":
				ctx.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, card)
}

type getCardRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getCard(ctx *gin.Context) {
	var req getCardRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	card, err := s.store.GetCard(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, card)
}

type listCardRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

func (s *Server) listCard(ctx *gin.Context) {
	var req listCardRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListMostRecentCardParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	card, err := s.store.ListMostRecentCard(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, card)
}
