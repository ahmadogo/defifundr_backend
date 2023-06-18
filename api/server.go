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
	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)
	router.POST("/user/verify", server.verifyUser)
	router.POST("/user/verify/resend", server.resendVerificationCode)
	router.POST("/user/password/reset", server.resetPassword)
	router.POST("/user/password/reset/verify", server.verifyPasswordResetCode)
	router.GET("/user/checkUsername/:username", server.checkUsernameExists)
	router.POST("/token/renewAccess", server.renewAccessToken)
	authRoutes := router.Group("/").Use(authMiddleWare(server.tokenMaker))
	authRoutes.GET("/user", server.getUser)
	authRoutes.GET("/user/avatar", server.setProfileAvatar)
	authRoutes.POST("/user/logout", server.logoutUser)
	authRoutes.POST("/user/password/change", server.changePassword)
	authRoutes.GET("/user/privatekey", server.getPrivateKey)
	authRoutes.GET("/campaigns", server.getCampaigns)
	authRoutes.POST("/campaigns", server.createCampaign)
	authRoutes.GET("/campaigns/:id", server.getCampaign)
	authRoutes.GET("/campaignsTypes", server.getCampaignTypes)
	authRoutes.GET("/campaigns/getDonors/:id", server.getCampaignDonors)
	authRoutes.POST("/campaigns/donate", server.donateToCampaign)
	authRoutes.POST("/campaigns/withdraw", server.withdrawFromCampaign)
	authRoutes.GET("/campaigns/myDonations", server.getMyDonations)
	authRoutes.GET("/campaigns/currentPrice", server.currentPrice)
	// authRoutes.GET("/campaigns/search", server.searchCampaigns)
	server.router = router
}

// Runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
