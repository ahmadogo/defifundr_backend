package api

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"

	db "github.com/demola234/defiraise/db/sqlc"
	crypt "github.com/demola234/defiraise/defi"
	"github.com/demola234/defiraise/interfaces"
	"github.com/demola234/defiraise/token"
	"github.com/demola234/defiraise/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get User
// @Description Get user details
// @Accept  json
// @Produce  json
// @Tags Profile
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success		200				{object}    interfaces.UserResponse{data=interfaces.UserResponse}	"success" 
// @Router /user [get]
func (server *Server) getUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload == nil {
		err := errors.New(interfaces.ErrUserNotFound)
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

	balance, err := crypt.GetBalance(user.Address)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	rsp := interfaces.UserResponse{
		Username:    user.Username,
		Email:       user.Email,
		Balance:     balance,
		Address:     user.Address,
		Biometrics:  user.Biometrics,
		Avatar:      user.Avatar,
		IsFirstTime: user.IsUsed,
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, rsp))
}

// @Summary Update User
// @Description Update user details
// @Accept  json
// @Produce  json
// @Tags Profile
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param   data        body   interfaces.CheckUsernameExistsRequest[types.Post]    true  "Get private key"
// @Success		200				string 	"User updated successfully"
// @Router /user/update [post]
func (server *Server) updateUser(ctx *gin.Context) {
	var req interfaces.CheckUsernameExistsRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload == nil {
		err := errors.New(interfaces.ErrUserNotFound)
		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
		return
	}

	_, err := server.store.UpdateUser(ctx, db.UpdateUserParams{
		Username: sql.NullString{
			String: authPayload.Username,
			Valid:  true,
		},
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, "User updated successfully"))
}

// @Summary Get User By Address
// @Description Get user details by address
// @Accept  json
// @Produce  json
// @Tags Profile
// @Param   data        body   interfaces.GetUserRequest[types.Post]    true  "Get private key"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success		200				{object}    interfaces.UserResponse{data=interfaces.UserResponse}	"success" 
// @Router /user/address [post]
func (server *Server) getUserByAddress(ctx *gin.Context) {
	var req interfaces.GetUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	user, err := server.store.GetUserByAddress(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	balance, err := crypt.GetBalance(user.Address)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	rsp := interfaces.UserResponse{
		Username:    user.Username,
		Email:       user.Email,
		Balance:     balance,
		Address:     user.Address,
		Biometrics:  user.Biometrics,
		Avatar:      user.Avatar,
		IsFirstTime: user.IsUsed,
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, rsp))
}

func (server *Server) logoutUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload == nil {
		err := errors.New(interfaces.ErrUserNotFound)
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

// @Summary Change Password
// @Description Change password of user
// @Accept  json
// @Produce  json
// @Tags Profile
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param   data        body   interfaces.ChangePasswordRequest[types.Post]    true  "Get private key"
// @Success		200				string 	"success"
// @Router /user/password/change [post]
func (server *Server) changePassword(ctx *gin.Context) {
	var req interfaces.ChangePasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload == nil {
		err := errors.New(interfaces.ErrUserNotFound)
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
		Username: sql.NullString{
			String: user.Username,
			Valid:  true,
		},
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

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, "Password changed successfully"))
}

// @Summary Get Private Key
// @Description Get private key of user
// @Accept  json
// @Produce  json
// @Tags Profile
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param   data        body   interfaces.GetPrivateKeyRequest[types.Post]    true  "Get private key"
// @Success		200				{object}    interfaces.UserResponse{data=interfaces.AddressResponse}	"success"
// @Router /user/privatekey [post]
func (server *Server) getPrivateKey(ctx *gin.Context) {
	var req interfaces.GetPrivateKeyRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload == nil {
		err := errors.New(interfaces.ErrUserNotFound)
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

	if err := utils.CheckPassword(req.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, interfaces.ErrorResponse(err, http.StatusUnauthorized))
		return
	}

	privateKey, address, err := crypt.DecryptPrivateKey(user.FilePath, server.config.PassPhase)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())

	arg := interfaces.AddressResponse{
		PrivateKey: privateKeyHex,
		Address:    address,
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, arg))
}

func (server *Server) setProfileAvatar(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload == nil {
		err := errors.New(interfaces.ErrUserNotFound)
		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
		return
	}

	imageURL, err := utils.UploadAvatar(ctx, file, authPayload.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	arg := db.UpdateUserParams{
		Username: sql.NullString{
			String: authPayload.Username,
			Valid:  true,
		},
		Avatar: sql.NullString{
			String: imageURL,
			Valid:  true,
		},
		IsFirstTime: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
	}

	_, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, "Image uploaded successfully"))
}

// @Summary Get Biometrics Settings
// @Description Get biometrics settings
// @Accept  json
// @Produce  json
// @Tags Profile
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param   data        body   interfaces.SetBiometricsRequest[types.Post]    true  "Get private key"
// @Success		200				{object}   string "Biometrics set successfully" 
// @Router /user/biometrics [post]
func (server *Server) setBiometrics(ctx *gin.Context) {
	var req interfaces.SetBiometricsRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload == nil {
		err := errors.New(interfaces.ErrUserNotFound)
		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
		return
	}

	arg := db.UpdateUserParams{
		Username: sql.NullString{
			String: authPayload.Username,
			Valid:  true,
		},
		Biometrics: sql.NullBool{
			Bool:  req.Biometrics,
			Valid: true,
		},
	}

	_, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, "Biometrics set successfully"))
}

func (server *Server) selectAvatar(ctx *gin.Context) {
	var req interfaces.Image
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload == nil {
		err := errors.New(interfaces.ErrUserNotFound)
		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
		return
	}

	imageURL := utils.GetAvatarUrl(req.ImageId - 1)

	arg := db.UpdateUserParams{
		Username: sql.NullString{
			String: authPayload.Username,
			Valid:  true,
		},
		Avatar: sql.NullString{
			String: imageURL,
			Valid:  true,
		},
		IsFirstTime: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
	}

	_, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
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

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, interfaces.NewUserResponse(user)))

}
