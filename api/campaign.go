package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/demola234/defiraise/crypto"
	crypt "github.com/demola234/defiraise/crypto"
	"github.com/demola234/defiraise/interfaces"
	"github.com/demola234/defiraise/token"
	"github.com/gin-gonic/gin"
)

func (server *Server) getCampaigns(ctx *gin.Context) {

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

	campaigns, err := crypto.GetCampaigns(user.Address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	camps := make([]interfaces.Campaigns, len(campaigns))

	for i, campaign := range campaigns {

		camps[i] = interfaces.Campaigns{
			CampaignType: campaign.CampaignType,
			Title:        campaign.Title,
			Deadline:     time.Now().Add(time.Duration(campaign.Deadline) * time.Second),
			Description:  campaign.Description,
			Goal:         campaign.Goal,
			Image:        campaign.Image,
		}
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, camps))
}

func (server *Server) getCampaign(ctx *gin.Context) {
	id := ctx.Param("id")
	// convert string id to int
	idL, err := strconv.Atoi(id)
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

	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	campaign, err := crypto.GetCampaign(idL, user.Address)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	// convert int to time.Time
	deadline := time.Unix(int64(campaign.Deadline), 0)

	camp := interfaces.Campaigns{
		CampaignType: campaign.CampaignType,
		Title:        campaign.Title,
		Description:  campaign.Description,
		Deadline:     deadline,
		Goal:         campaign.Goal,
		Image:        campaign.Image,
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, camp))
}

func (server *Server) getCampaignTypes(ctx *gin.Context) {
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

	campaignTypes, err := crypto.GetCampaignTypes(user.Address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, campaignTypes))
}

func (server *Server) getCampaignDonors(ctx *gin.Context) {
	id := ctx.Param("id")
	// convert string id to int
	idL, err := strconv.Atoi(id)
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

	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	donators, donations, totalFunds, err := crypto.GetDonorsAddressesAndAmounts(idL, user.Address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	for i, donor := range donators {
		donators[i] = donor
	}

	for i, donation := range donations {
		donations[i] = donation
	}

	// common.Address to string
	addresses := make([]string, len(donators))

	donationsList := make([]string, len(donations))

	donors := interfaces.Donations{
		Address:   addresses,
		Amount:    totalFunds.Int64(),
		Donations: donationsList,
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, donors))
}

func (server *Server) donateToCampaign(ctx *gin.Context) {
	var donation interfaces.Donation

	err := ctx.ShouldBindJSON(&donation)
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

	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	privateKey, address, err := crypt.DecryptPrivateKey(user.FilePath, server.config.PassPhase)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	msg, err := crypto.Donate(donation.Amount, donation.CampaignId, privateKey, address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, msg))
}

func (server *Server) createCampaign(ctx *gin.Context) {
	var campaign interfaces.Campaigns

	err := ctx.ShouldBindJSON(&campaign)
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

	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	privateKey, address, err := crypt.DecryptPrivateKey(user.FilePath, server.config.PassPhase)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}
	// int64 to int
	campaignGoal := int(campaign.Goal)
	// add deadline to current time
	dl := time.Unix(int64(campaign.Deadline.Day()), 0)

	campaigns, err := crypto.CreateCampaign(campaign.Title, campaign.CampaignType, campaign.Description, campaignGoal, dl, campaign.Image, privateKey, address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, campaigns))
}

func (server *Server) getMyDonations(ctx *gin.Context) {
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

	donations, err := crypto.GetCampaignsByOwner(user.Address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, donations))
}

func (server *Server) withdrawFromCampaign(ctx *gin.Context) {
	var withdraw interfaces.Withdraw

	err := ctx.ShouldBindJSON(&withdraw)
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

	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	privateKey, address, err := crypt.DecryptPrivateKey(user.FilePath, server.config.PassPhase)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	msg, err := crypto.PayOut(withdraw.CampaignId, address, privateKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, msg))
}

func (server *Server) currentEthPrice(ctx *gin.Context) {
	price, err := crypto.GetEthPrice()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, price))
}
