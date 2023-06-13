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

	username := req.FirstName

	arg := db.CreateUserParams{
		Username:       username,
		HashedPassword: req.Password,
		FirstName:      req.FirstName,
		Email:          req.Email,
		Avatar:         "",
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	rsp := interfaces.NewUserResponse(user)
	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, rsp))

}
