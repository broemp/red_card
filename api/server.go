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
	if gin.Mode() == gin.DebugMode {
		router.Use(cors.Default())
		return
	}
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.CORS_Frontend},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
}
