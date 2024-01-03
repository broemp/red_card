package api

import (
	"fmt"
	"time"

	db "github.com/broemp/red_card/db/sqlc"
	"github.com/broemp/red_card/token"
	"github.com/broemp/red_card/util"
	"github.com/gin-contrib/cors"
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
	addCors(router, s.config)
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

func addCors(router *gin.Engine, config util.Config) {
	if config.CORS_Enable {
		if config.CORS_Frontend == "" {
			fmt.Println("Failed to read Allowed Origin. Please Provide Origion via CORS_Frontend!")
			router.Use(cors.Default())
		}

		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{config.CORS_Frontend},
			AllowMethods:     []string{"*"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"*"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
		return
	}
}
