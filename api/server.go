package api

import (
	db "github.com/demola234/defiraise/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {

	server := &Server{
		store:  store,
		router: gin.Default(),
	}

	server.setUpRouter()
	return server
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	server.router = router
}

// Runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
