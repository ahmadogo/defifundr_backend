package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/demola234/defiraise/crypto"
	db "github.com/demola234/defiraise/db/sqlc"
	"github.com/demola234/defiraise/interfaces"
	"github.com/demola234/defiraise/token"
	"github.com/demola234/defiraise/utils"
	"github.com/gin-gonic/gin"
)

func (server *Server) getUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload == nil {
		err := errors.New("user not found")
		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
		return
	}

	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	balance, err := crypto.GetBalance(user.Address)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	rsp := interfaces.UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Balance:  balance,
		Address:  user.Address,
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, rsp))
}

func (server *Server) logoutUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload == nil {
		err := errors.New("user not found")
		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
		return
	}

	_, err := server.store.DeleteSession(ctx, authPayload.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, nil))
}

func (server *Server) changePassword(ctx *gin.Context) {
	var req interfaces.ChangePasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload == nil {
		err := errors.New("user not found")
		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
		return
	}

	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	if err := utils.CheckPassword(req.OldPassword, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, interfaces.ErrorResponse(err, http.StatusUnauthorized))
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	arg := db.UpdateUserParams{
		Username: user.Username,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
	}

	_, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, nil))
}
