package api

import (
	"net/http"

	db "github.com/demola234/defiraise/db/sqlc"
	"github.com/demola234/defiraise/interfaces"
	"github.com/gin-gonic/gin"
)

func (server *Server) createUser(ctx *gin.Context) {
	var req interfaces.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	username := req.Username

	arg := db.CreateUserParams{
		Username:       username,
		HashedPassword: req.Password,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	// filepath, address, err := crypto.GenerateAccountKeyStone(password)

	rsp := interfaces.NewUserResponse(user)
	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, rsp))

}
