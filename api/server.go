package api

import (
	"fmt"

	db "github.com/broemp/red_card/db/sqlc"
	"github.com/broemp/red_card/token"
	"github.com/broemp/red_card/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.JWT_SECRET)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()

	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()
	router.Use(corsMiddleware())
	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker))

	router.POST("/users/register", s.registerUser)
	router.POST("/users/login", s.loginUser)
	router.GET("/users/:id", s.getUser)
	router.GET("/users/:id/cards", s.getUserCards)
	authRoutes.GET("/users/", s.listUserFilter)

	router.GET("/cards", s.listCard)
	router.GET("/cards/:id", s.getCard)

	authRoutes.POST("/cards", s.createCard)

	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(":" + address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
