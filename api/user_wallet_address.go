package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	db "github.com/demola234/defiraise/db/sqlc"
	"github.com/demola234/defiraise/interfaces"
	"github.com/demola234/defiraise/token"
	"github.com/gin-gonic/gin"
)

// @Summary Create user Wallet Address
// @Description Create a new user wallet address
// @Accept  json
// @Produce  json
// @Tags User Wallet Addresses
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body interfaces.CreateWalletAddressRequest true "Wallet address Creation Data"
// @Success 200 {object} interfaces.DocSuccessResponse{data=interfaces.AddressResponse} "success"
// @Router /wallet-address/create [post]
func (server *Server) createWalletAddress(ctx *gin.Context) {
	var req interfaces.CreateWalletAddressRequest

	// Validate request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(errors.New("invalid request"), http.StatusBadRequest))
		return
	}

	// Extract authenticated user from context
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload == nil {
		ctx.JSON(http.StatusUnauthorized, interfaces.ErrorResponse(errors.New("unauthorized"), http.StatusUnauthorized))
		return
	}

	// Ensure the user exists
	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(errors.New("user not found"), http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	// Create the wallet entry
	wallet, err := server.store.CreateUserWallet(ctx, db.CreateUserWalletParams{
		UserID:        user.Username,
		WalletAddress: req.WalletAddress,
		Chain:         req.Chain,
		Status:        db.UserWalletAddressesStatusesActive, // Default status: active
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	// Response
	response := interfaces.WalletAddressResponse{
		ID:            wallet.ID,
		UserID:        wallet.UserID,
		WalletAddress: wallet.WalletAddress,
		Chain:         wallet.Chain,
		Status:        string(wallet.Status),
		CreatedAt:     wallet.CreatedAt,
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, response))
}

// @Summary Get User Wallets (Paginated)
// @Description Fetch paginated wallets for the authenticated user
// @Accept  json
// @Produce  json
// @Tags User Wallet Addresses
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param limit query int false "Limit (default: 10)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} interfaces.DocSuccessResponse{data=[]interfaces.WalletAddressResponse} "success"
// @Router /wallet-address [get]
func (server *Server) getUserWallets(ctx *gin.Context) {
	// Extract pagination params
	limit := ctx.DefaultQuery("limit", "10")  // Default: 10 items
	offset := ctx.DefaultQuery("offset", "0") // Default: start from 0

	// Convert limit and offset to integers
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(errors.New("invalid limit value"), http.StatusBadRequest))
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt < 0 {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(errors.New("invalid offset value"), http.StatusBadRequest))
		return
	}

	// Extract authenticated user
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload == nil {
		ctx.JSON(http.StatusUnauthorized, interfaces.ErrorResponse(errors.New("unauthorized"), http.StatusUnauthorized))
		return
	}

	// Ensure user exists
	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(errors.New("user not found"), http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	var getUserWalletsParams = db.GetUserWalletsParams{
		UserID: user.Username,
		Limit:  int32(limitInt),
		Offset: int32(offsetInt),
	}

	// Fetch paginated wallets
	wallets, err := server.store.GetUserWallets(ctx, getUserWalletsParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	// Transform wallets into response format
	walletResponses := make([]interfaces.WalletAddressResponse, len(wallets))
	for i, wallet := range wallets {
		walletResponses[i] = interfaces.WalletAddressResponse{
			ID:            wallet.ID,
			UserID:        wallet.UserID,
			WalletAddress: wallet.WalletAddress,
			Chain:         wallet.Chain,
			Status:        string(wallet.Status),
			CreatedAt:     wallet.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, walletResponses))
}


// @Summary Get Wallet by ID
// @Description Fetch wallet if it belongs to the authenticated user
// @Accept json
// @Produce json
// @Tags User Wallet Addresses
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id path int64 true "Wallet ID"
// @Success 200 {object} interfaces.AddressResponse "success"
// @Router /wallet-address/{id} [get]
func (server *Server) getWalletByID(ctx *gin.Context) {
	var req struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(errors.New("invalid ID"), http.StatusBadRequest))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload == nil {
		ctx.JSON(http.StatusUnauthorized, interfaces.ErrorResponse(errors.New("unauthorized"), http.StatusUnauthorized))
		return
	}

	// Ensure the user exists
	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(errors.New("user not found"), http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	var getWalletByIDParams = db.GetWalletByIdParams {
		ID: req.ID,
		UserID: user.Username,
	}

	// Fetch wallet and verify ownership
	wallet, err := server.store.GetWalletById(ctx, getWalletByIDParams)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(errors.New("wallet not found"), http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	if wallet.UserID != authPayload.Username {
		ctx.JSON(http.StatusForbidden, interfaces.ErrorResponse(errors.New("not authorized"), http.StatusForbidden))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, wallet))
}

// @Summary Get Wallet by Address
// @Description Fetch wallet by address
// @Accept json
// @Produce json
// @Tags User Wallet Addresses
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param wallet_address path string true "Wallet Address"
// @Success 200 {object} interfaces.AddressResponse "success"
// @Router /wallet-address/address/{wallet_address} [get]
func (server *Server) getWalletByAddress(ctx *gin.Context) {
	var req struct {
		WalletAddress string `uri:"wallet_address" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(errors.New("invalid request"), http.StatusBadRequest))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload == nil {
		ctx.JSON(http.StatusUnauthorized, interfaces.ErrorResponse(errors.New("unauthorized"), http.StatusUnauthorized))
		return
	}

	// Ensure the user exists
	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(errors.New("user not found"), http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	var getWalletByAddressParams = db.GetWalletByAddressParams {
		WalletAddress: req.WalletAddress,
		UserID: user.Username,
	}

	wallet, err := server.store.GetWalletByAddress(ctx, getWalletByAddressParams)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(errors.New("wallet not found"), http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, wallet))
}

// @Summary Update Wallet Status
// @Description Update the status of a user's wallet
// @Accept json
// @Produce json
// @Tags User Wallet Addresses
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body interfaces.UpdateUserWalletAddressStatusRequest true "Wallet status update data"
// @Success 200 {object} interfaces.AddressResponse "success"
// @Router /wallet-address/status [patch]
func (server *Server) updateWalletStatus(ctx *gin.Context) {
	var req interfaces.UpdateUserWalletAddressStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(errors.New("invalid request"), http.StatusBadRequest))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload == nil {
		ctx.JSON(http.StatusUnauthorized, interfaces.ErrorResponse(errors.New("unauthorized"), http.StatusUnauthorized))
		return
	}

	// Ensure the user exists
	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(errors.New("user not found"), http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	var getWalletByIdParams = db.GetWalletByIdParams {
		ID: req.ID,
		UserID: user.Username,
	}

	wallet, err := server.store.GetWalletById(ctx, getWalletByIdParams)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(errors.New("wallet not found"), http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	if wallet.UserID != authPayload.Username {
		ctx.JSON(http.StatusForbidden, interfaces.ErrorResponse(errors.New("not authorized"), http.StatusForbidden))
		return
	}

	var updateUserWalletStatusParams = db.UpdateUserWalletStatusParams {
		ID: req.ID,
		UserID: user.Username,
		Status:        db.UserWalletAddressesStatuses(req.Status),
	}

	fmt.Println(updateUserWalletStatusParams);


	updatedWallet, err := server.store.UpdateUserWalletStatus(ctx, updateUserWalletStatusParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	fmt.Println(updatedWallet);

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, updatedWallet))
}

// @Summary Soft Delete Wallet
// @Description Soft delete a wallet (marks as deleted and updates status)
// @Accept json
// @Produce json
// @Tags User Wallet Addresses
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id path string true "Wallet Address Id"
// @Success 200 {object} interfaces.AddressResponse "success"
// @Router /wallet-address/{id} [delete]
func (server *Server) softDeleteWallet(ctx *gin.Context) {
	var req struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(errors.New("invalid request"), http.StatusBadRequest))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload == nil {
		ctx.JSON(http.StatusUnauthorized, interfaces.ErrorResponse(errors.New("unauthorized"), http.StatusUnauthorized))
		return
	}

	// Ensure the user exists
	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(errors.New("user not found"), http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	var getWalletByIdParams = db.GetWalletByIdParams {
		ID: req.ID,
		UserID: user.Username,
	}

	wallet, err := server.store.GetWalletById(ctx, getWalletByIdParams)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(errors.New("wallet not found"), http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	if wallet.UserID != authPayload.Username {
		ctx.JSON(http.StatusForbidden, interfaces.ErrorResponse(errors.New("not authorized"), http.StatusForbidden))
		return
	}

	var softDeleteUserWalletParams = db.SoftDeleteUserWalletParams {
		ID: req.ID,
		UserID: user.Username,
	}

	updatedWallet, err := server.store.SoftDeleteUserWallet(ctx, softDeleteUserWalletParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, updatedWallet))
}