package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/demola234/defiraise/crypto"
	db "github.com/demola234/defiraise/db/sqlc"
	"github.com/demola234/defiraise/interfaces"
	"github.com/demola234/defiraise/utils"
	"github.com/gin-gonic/gin"
)

func (server *Server) createUser(ctx *gin.Context) {
	var req interfaces.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	configs, err := utils.LoadConfig("./../")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	filepath, address, err := crypto.GenerateAccountKeyStone(configs.PassPhase)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	email := utils.EmailInfo{
		Name:    req.Username,
		Details: "Welcome to DefiRaise",
		Otp:     utils.RandomOtp(),
		Subject: "Welcome to DefiRaise",
	}

	// hash password
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	arg := db.CreateUserParams{
		Username:        req.Username,
		HashedPassword:  hashPassword,
		Email:           req.Email,
		Address:         address,
		FilePath:        filepath,
		Balance:         "0",
		IsUsed:          false,
		SecretCode:      email.Otp,
		IsEmailVerified: false,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, interfaces.ErrorResponse(err, http.StatusForbidden))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	_, err = utils.SendEmail(req.Email, req.Username, email, "./utils")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	rsp := interfaces.NewUserResponse(user)
	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, rsp))

}

func (server *Server) loginUser(ctx *gin.Context) {
	var req interfaces.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	if !user.IsEmailVerified {
		err := errors.New("user not verified")
		ctx.JSON(http.StatusForbidden, interfaces.ErrorResponse(err, http.StatusForbidden))
		return
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, interfaces.ErrorResponse(err, http.StatusUnauthorized))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	// get current eth balance
	balance, err := crypto.GetBalance(user.Address)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	// update user balance
	user, err = server.store.UpdateUser(ctx, db.UpdateUserParams{
		Username: user.Username,
		Balance: sql.NullString{
			String: balance,
			Valid:  true,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	sessions, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		Username:     user.Username,
		RefreshToken: refreshToken,
		ID:           refreshPayload.ID,
		ExpiresAt:    refreshPayload.ExpiresAt,
		UserAgent:    ctx.GetHeader("User-Agent"),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	rsp := interfaces.LoginResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		SessionID:             sessions.ID,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt,
		User:                  interfaces.NewUserResponse(user),
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, rsp))
}

func (server *Server) verifyUser(ctx *gin.Context) {
	var req interfaces.VerifyUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
		return
	}

	if user.IsEmailVerified {
		err := errors.New("user already verified")
		ctx.JSON(http.StatusForbidden, interfaces.ErrorResponse(err, http.StatusForbidden))
		return
	}

	if user.IsUsed {
		err := errors.New("user already used")
		ctx.JSON(http.StatusForbidden, interfaces.ErrorResponse(err, http.StatusForbidden))
		return
	}

	if user.ExpiredAt.Before(time.Now()) {
		err := errors.New("otp has expired")
		ctx.JSON(http.StatusForbidden, interfaces.ErrorResponse(err, http.StatusForbidden))
		return
	}

	if user.SecretCode != req.OtpCode {
		err := errors.New("invalid otp code")
		ctx.JSON(http.StatusForbidden, interfaces.ErrorResponse(err, http.StatusForbidden))
		return
	}

	arg := db.UpdateUserParams{
		Username: user.Username,
		IsEmailVerified: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		IsUsed: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
	}

	_, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, "User Verified"))
}

func (server *Server) resendVerificationCode(ctx *gin.Context) {
	var req interfaces.ResendVerificationCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
		return
	}

	if user.IsEmailVerified {
		err := errors.New("user already verified")
		ctx.JSON(http.StatusForbidden, interfaces.ErrorResponse(err, http.StatusForbidden))
		return
	}

	if user.IsUsed {
		err := errors.New("user already used")
		ctx.JSON(http.StatusForbidden, interfaces.ErrorResponse(err, http.StatusForbidden))
		return
	}

	email := utils.EmailInfo{
		Name:    user.Username,
		Details: "Welcome to DefiRaise",
		Otp:     utils.RandomOtp(),
		Subject: "Welcome to DefiRaise",
	}

	arg := db.UpdateUserParams{
		Username: user.Username,
		SecretCode: sql.NullString{
			String: email.Otp,
			Valid:  true,
		},
		ExpiredAt: sql.NullTime{
			Time:  time.Now().Add(15 * time.Minute),
			Valid: true,
		},
	}

	_, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	_, err = utils.SendEmail(user.Email, user.Username, email, "./utils")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, "OTP code resent"))
}
