package api

import (
	"fmt"

	db "github.com/demola234/defiraise/db/sqlc"
	"github.com/demola234/defiraise/token"
	"github.com/demola234/defiraise/utils"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewTokenMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %s", err.Error())
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		router:     gin.Default(),
	}

	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	authRoutes := router.Group("/").Use(authMiddleWare(server.tokenMaker))
	authRoutes.GET("/users", server.getUser)
	// authRoutes.PATCH("/users/changePassword", server.changePassword)

	server.router = router
}

// Runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
