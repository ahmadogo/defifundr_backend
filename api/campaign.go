package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/demola234/defiraise/defi"
	crypt "github.com/demola234/defiraise/defi"
	"github.com/demola234/defiraise/interfaces"
	"github.com/demola234/defiraise/token"
	"github.com/demola234/defiraise/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get campaigns
// @Description Get campaigns
// @Accept  json
// @Produce  json
// @Tags Campaigns
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success		200				{object}    interfaces.DocSuccessResponse{data=[]interfaces.Campaigns}	"success"
// @Router /campaigns [get]
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

	campaigns, err := defi.GetCampaigns(user.Address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	camps := make([]interfaces.Campaigns, len(campaigns))

	// if camps is empty
	if len(camps) == 0 {
		ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, []string{}))
		return
	}

	for i, campaign := range campaigns {
		userInfo, _ := server.store.GetUserByAddress(ctx, campaign.Owner)

		totalNumber, err := defi.GetTotalDonationsByCampaignId(i)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
			return
		}

		//  Loop through the donators and amounts
		donators, amounts, _, err := defi.GetDonorsAddressesAndAmounts(i)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
			return
		}

		dons := make([]interfaces.DonorDetails, len(donators))

		if len(donators) != 0 {
			for k, _ := range donators {
				getUser, _ := server.store.GetUserByAddress(ctx, donators[k])

				dons[k] = interfaces.DonorDetails{
					Amount:   (float64(amounts[k]) / 1e18),
					Donor:    donators[k],
					Image:    getUser.Avatar,
					Username: getUser.Username,
				}

				// wei to either

				// if the deadline has been reached and skipped remove the campaign from the list of campaigns

				camps[i] = interfaces.Campaigns{
					CampaignType:       campaign.CampaignType,
					Title:              campaign.Title,
					Deadline:           time.Unix(int64(campaign.Deadline), 0),
					Description:        campaign.Description,
					Goal:               float64(campaign.Goal),
					Image:              campaign.Image,
					TotalAmountDonated: float64(campaign.TotalFunds),
					TotalNumber:        totalNumber.Int64(),
					Owner:              campaign.Owner,
					ID:                 int(campaign.ID),
					Donations:          dons,
					User: []interfaces.UserResponseInfo{
						{
							Username: userInfo.Username,
							Email:    userInfo.Email,
							Address:  userInfo.Address,
							Avatar:   userInfo.Avatar,
						},
					},
				}
			}
		} else {
			// skip the current iteration and remove empty campaign with empty description and title from list of campaigns to be displayed to the user on the frontend side of the application

			// remove campaign if deadline is less than current time do not display to the user
			// if deadline is less than current time

			camps[i] = interfaces.Campaigns{
				CampaignType:       campaign.CampaignType,
				Title:              campaign.Title,
				Deadline:           time.Unix(int64(campaign.Deadline), 0),
				Description:        campaign.Description,
				Goal:               float64(campaign.Goal),
				Image:              campaign.Image,
				TotalAmountDonated: float64(campaign.TotalFunds),
				TotalNumber:        totalNumber.Int64(),
				Owner:              campaign.Owner,
				ID:                 int(campaign.ID),
				Donations:          dons,
				User: []interfaces.UserResponseInfo{
					{
						Username: userInfo.Username,
						Email:    userInfo.Email,
						Address:  userInfo.Address,
						Avatar:   userInfo.Avatar,
					},
				},
			}

		}
	}

	//  if deadline is less than current time remove the element from the list of campaigns to be displayed to the user

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, camps))
}

// @Summary Get latest active campaigns
// @Description Get latest active campaigns
// @Accept  json
// @Produce  json
// @Tags Campaigns
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success		200				{object}    interfaces.DocSuccessResponse{data=[]interfaces.Campaigns}	"success"
// @Router /campaigns/latestCampaigns [get]
func (server *Server) getLatestActiveCampaigns(ctx *gin.Context) {
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

	campaigns, err := defi.GetCampaigns(user.Address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	activeCampaigns := []interfaces.Campaigns{} // Slice to store active campaigns

	for i, campaign := range campaigns {
		deadline := time.Unix(int64(campaign.Deadline), 0)
		userInfo, _ := server.store.GetUserByAddress(ctx, campaign.Owner)

		totalNumber, err := defi.GetTotalDonationsByCampaignId(i)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
			return
		}

		//  Loop through the donators and amounts
		donators, amounts, _, err := defi.GetDonorsAddressesAndAmounts(i)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
			return
		}

		dons := make([]interfaces.DonorDetails, len(donators))

		if deadline.After(time.Now()) {
			if len(donators) != 0 {
				for k, _ := range donators {
					getUser, _ := server.store.GetUserByAddress(ctx, donators[k])

					dons[k] = interfaces.DonorDetails{
						Amount:   (float64(amounts[k]) / 1e18),
						Donor:    donators[k],
						Image:    getUser.Avatar,
						Username: getUser.Username,
					}

					activeCampaigns = append(activeCampaigns, interfaces.Campaigns{
						CampaignType:       campaign.CampaignType,
						Title:              campaign.Title,
						Deadline:           deadline,
						Description:        campaign.Description,
						Goal:               float64(campaign.Goal),
						Image:              campaign.Image,
						TotalAmountDonated: float64(campaign.TotalFunds),
						TotalNumber:        totalNumber.Int64(),
						Owner:              campaign.Owner,
						ID:                 int(campaign.ID),
						Donations:          dons,
						User: []interfaces.UserResponseInfo{
							{
								Username: userInfo.Username,
								Email:    userInfo.Email,
								Address:  userInfo.Address,
								Avatar:   userInfo.Avatar,
							},
						},
					})

				}
			} else {
				activeCampaigns = append(activeCampaigns, interfaces.Campaigns{
					CampaignType:       campaign.CampaignType,
					Title:              campaign.Title,
					Deadline:           deadline,
					Description:        campaign.Description,
					Goal:               float64(campaign.Goal),
					Image:              campaign.Image,
					TotalAmountDonated: float64(campaign.TotalFunds),
					TotalNumber:        totalNumber.Int64(),
					Owner:              campaign.Owner,
					ID:                 int(campaign.ID),
					Donations:          dons,
					User: []interfaces.UserResponseInfo{
						{
							Username: userInfo.Username,
							Email:    userInfo.Email,
							Address:  userInfo.Address,
							Avatar:   userInfo.Avatar,
						},
					},
				})
			}
		}
	}

	// Sort the active campaigns by deadline in ascending order
	for i := 0; i < len(activeCampaigns); i++ {
		for j := i + 1; j < len(activeCampaigns); j++ {
			if activeCampaigns[i].Deadline.After(activeCampaigns[j].Deadline) {
				temp := activeCampaigns[i]
				activeCampaigns[i] = activeCampaigns[j]
				activeCampaigns[j] = temp
			}
		}
	}
	latestActiveCampaigns := []interfaces.Campaigns{}

	// Get the latest 10 campaigns or all active campaigns if less than 10
	if len(activeCampaigns) <= 10 {
		latestActiveCampaigns = activeCampaigns
	} else {
		latestActiveCampaigns = activeCampaigns[:10]
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, latestActiveCampaigns))
}

// @Summary Get Campaigns by category
// @Description Get Campaigns by category
// @Accept  json
// @Produce  json
// @Tags Campaigns
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id path int true "Category ID"
// @Success		200				{object}    interfaces.DocSuccessResponse{data=[]interfaces.Campaigns}	"success"
// @Router /campaigns/categories/{id} [get]
func (server *Server) getCampaignsByCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	// convert string id to int
	idL, err := strconv.Atoi(id)
	if err != nil {
		errors := errors.New("invalid id")
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(errors, http.StatusBadRequest))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload == nil {
		err := errors.New(interfaces.ErrUserNotFound)
		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
		return
	}

	// _, err := server.store.GetUser(ctx, authPayload.Username)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
	// 	return
	// }

	campaigns, err := defi.GetCampaignByCategory(int64(idL))

	fmt.Println(campaigns)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	activeCampaigns := []interfaces.Campaigns{} // Slice to store active ca

	// if camps is empty
	if len(activeCampaigns) == 0 {
		ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, []string{}))
		return
	}

	for i, campaign := range campaigns {
		// skip the current iteration and remove empty campaign with empty description and title from list of campaigns to be displayed to the user on the frontend side of the application

		userInfo, _ := server.store.GetUserByAddress(ctx, campaign.Owner)

		totalNumber, err := defi.GetTotalDonationsByCampaignId(i)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
			return
		}

		//  Loop through the donators and amounts
		donators, amounts, _, err := defi.GetDonorsAddressesAndAmounts(i)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
			return
		}

		dons := make([]interfaces.DonorDetails, len(donators))

		// convert unix to time.Time
		// deadline := time.Unix(int64(campaign.Deadline), 0)

		if len(donators) != 0 {
			for k, _ := range donators {
				getUser, _ := server.store.GetUserByAddress(ctx, donators[k])

				dons[k] = interfaces.DonorDetails{
					Amount:   (float64(amounts[k]) / 1e18),
					Donor:    donators[k],
					Image:    getUser.Avatar,
					Username: getUser.Username,
				}
			}
			// wei to either

			// if the deadline has been reached and skipped remove the campaign from the list of campaigns

			activeCampaigns = append(activeCampaigns, interfaces.Campaigns{
				CampaignType:       campaign.CampaignType,
				Title:              campaign.Title,
				Deadline:           time.Unix(int64(campaign.Deadline), 0),
				Description:        campaign.Description,
				Goal:               float64(campaign.Goal),
				Image:              campaign.Image,
				TotalAmountDonated: float64(campaign.TotalFunds),
				TotalNumber:        totalNumber.Int64(),
				Owner:              campaign.Owner,
				ID:                 int(campaign.ID),
				Donations:          dons,
				User: []interfaces.UserResponseInfo{
					{
						Username: userInfo.Username,
						Email:    userInfo.Email,
						Address:  userInfo.Address,
						Avatar:   userInfo.Avatar,
					},
				},
			})

		} else {
			// skip the current iteration and remove empty campaign with empty description and title from list of campaigns to be displayed to the user on the frontend side of the application

			// remove campaign if deadline is less than current time do not display to the user
			// if deadline is less than current time

			activeCampaigns = append(activeCampaigns, interfaces.Campaigns{
				CampaignType:       campaign.CampaignType,
				Title:              campaign.Title,
				Deadline:           time.Unix(int64(campaign.Deadline), 0),
				Description:        campaign.Description,
				Goal:               float64(campaign.Goal),
				Image:              campaign.Image,
				TotalAmountDonated: float64(campaign.TotalFunds),
				TotalNumber:        totalNumber.Int64(),
				Owner:              campaign.Owner,
				ID:                 int(campaign.ID),
				Donations:          dons,
				User: []interfaces.UserResponseInfo{
					{
						Username: userInfo.Username,
						Email:    userInfo.Email,
						Address:  userInfo.Address,
						Avatar:   userInfo.Avatar,
					},
				},
			})
		}
	}

	//  if deadline is less than current time remove the element from the list of campaigns to be displayed to the user

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, activeCampaigns))

}

// @Summary Get Campaigns by owner
// @Description Get Campaigns by owner
// @Accept  json
// @Produce  json
// @Tags Campaigns
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success		200				{object}    interfaces.DocSuccessResponse{data=[]interfaces.Campaigns}	"success"
// @Router /campaigns/owner [get]
func (server *Server) getCampaignsByOwner(ctx *gin.Context) {
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

	campaigns, err := defi.GetCampaignsByOwner(user.Address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	activeCampaigns := []interfaces.Campaigns{} // Slice to store active campaigns

	for i, campaign := range campaigns {
		deadline := time.Unix(int64(campaign.Deadline), 0)
		userInfo, _ := server.store.GetUserByAddress(ctx, campaign.Owner)

		totalNumber, err := defi.GetTotalDonationsByCampaignId(i)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
			return
		}

		//  Loop through the donators and amounts
		donators, amounts, _, err := defi.GetDonorsAddressesAndAmounts(i)
		// fmt.Println(donators)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
			return
		}

		// Loop through the donators and amounts
		// map the donators and amounts to the DonorDetails struct and append to the dons slice of DonorDetails if the index is same as the index of the campaign

		dons := make([]interfaces.DonorDetails, len(donators))

		if deadline.After(time.Now()) {
			if len(donators) != 0 {
				for k, _ := range donators {
					getUser, _ := server.store.GetUserByAddress(ctx, donators[k])

					dons[k] = interfaces.DonorDetails{
						Amount:   (float64(amounts[k]) / 1e18),
						Donor:    donators[k],
						Image:    getUser.Avatar,
						Username: getUser.Username,
					}
				}

				activeCampaigns = append(activeCampaigns, interfaces.Campaigns{
					CampaignType:       campaign.CampaignType,
					Title:              campaign.Title,
					Deadline:           deadline,
					Description:        campaign.Description,
					Goal:               float64(campaign.Goal),
					Image:              campaign.Image,
					TotalAmountDonated: float64(campaign.TotalFunds),
					TotalNumber:        totalNumber.Int64(),
					Owner:              campaign.Owner,
					ID:                 int(campaign.ID),
					Donations:          dons,
					User: []interfaces.UserResponseInfo{
						{
							Username: userInfo.Username,
							Email:    userInfo.Email,
							Address:  userInfo.Address,
							Avatar:   userInfo.Avatar,
						},
					},
				})

			} else {
				activeCampaigns = append(activeCampaigns, interfaces.Campaigns{
					CampaignType:       campaign.CampaignType,
					Title:              campaign.Title,
					Deadline:           deadline,
					Description:        campaign.Description,
					Goal:               float64(campaign.Goal),
					Image:              campaign.Image,
					TotalAmountDonated: float64(campaign.TotalFunds),
					TotalNumber:        totalNumber.Int64(),
					Owner:              campaign.Owner,
					ID:                 int(campaign.ID),
					Donations:          dons,
					User: []interfaces.UserResponseInfo{
						{
							Username: userInfo.Username,
							Email:    userInfo.Email,
							Address:  userInfo.Address,
							Avatar:   userInfo.Avatar,
						},
					},
				})
			}
		}
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, activeCampaigns))
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

	if user.Address == "" {
		err := errors.New(interfaces.ErrUserNotFound)
		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
		return
	}

	campaign, err := defi.GetCampaign(idL)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	totalNumber, err := defi.GetTotalDonationsByCampaignId(idL)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	donators, amounts, _, err := defi.GetDonorsAddressesAndAmounts(idL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	dons := make([]interfaces.DonorDetails, len(donators))

	if len(donators) != 0 {
		for k, _ := range donators {
			getUser, _ := server.store.GetUserByAddress(ctx, donators[k])

			dons[k] = interfaces.DonorDetails{
				Amount:   (float64(amounts[k]) / 1e18),
				Donor:    donators[k],
				Image:    getUser.Avatar,
				Username: getUser.Username,
			}
		}

	}

	// convert int to time.Time
	deadline := time.Unix(int64(campaign.Deadline), 0)

	camp := interfaces.Campaigns{
		CampaignType: campaign.CampaignType,
		Title:        campaign.Title,
		Description:  campaign.Description,
		Deadline:     deadline,
		Goal:         float64(campaign.Goal / 1000000000000000000),
		Image:        campaign.Image,
		TotalNumber:  totalNumber.Int64(),
		Owner:        campaign.Owner,
		ID:           int(campaign.ID),
		User: []interfaces.UserResponseInfo{
			{
				Username: user.Username,
				Email:    user.Email,
				Address:  user.Address,
				Avatar:   user.Avatar,
			},
		},
		TotalAmountDonated: float64(campaign.TotalFunds / 1000000000000000000),
		Donations:          dons,
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

	campaignTypes, err := defi.GetCampaignTypes(user.Address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, campaignTypes))
}

// @Summary Get Campaign Donors
// @Description Get Campaign Donors
// @Accept  json
// @Produce  json
// @Tags Campaigns
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id path int true "Campaign ID"
// @Success		200				{object}    interfaces.DocSuccessResponse{data=[]interfaces.DonorDetails}	"success"
// @Router /campaigns/donors/{id} [get]
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

	// user, err := server.store.GetUser(ctx, authPayload.Username)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
	// 	return
	// }

	donators, amounts, _, err := defi.GetDonorsAddressesAndAmounts(idL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	dons := make([]interfaces.DonorDetails, len(donators))

	if len(donators) != 0 {
		for k, _ := range donators {
			getUser, _ := server.store.GetUserByAddress(ctx, donators[k])

			dons[k] = interfaces.DonorDetails{
				Amount:   (float64(amounts[k]) / 1e18),
				Donor:    donators[k],
				Image:    getUser.Avatar,
				Username: getUser.Username,
			}
		}

	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, dons))
}

// @Summary Donate to campaign
// @Description Donate to campaign
// @Accept  json
// @Produce  json
// @Tags Campaigns
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param   data        body   interfaces.Donation[types.Post]    true  "Donation"
// @Success		200				{object}    interfaces.DocSuccessResponse	"success"
// @Router /campaigns/donate [post]
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
		newErr := errors.New("unable to make transaction at this time, please try again later")
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(newErr, http.StatusBadRequest))
		return
	}

	// convert string id to int
	idL, err := strconv.Atoi(donation.CampaignId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	// convert string to float64
	amount, err := strconv.ParseFloat(donation.Amount, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	// convert balance from string to float64
	balance, err := strconv.ParseFloat(user.Balance, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	campaign, err := defi.GetCampaign(idL)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	// check if campaign is still active and not expired
	deadline := time.Unix(int64(campaign.Deadline), 0)
	if time.Now().After(deadline) {
		newErr := errors.New("campaign has closed")
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(newErr, http.StatusBadRequest))
		return
	}

	// check if user has enough balance
	if float64(amount) > balance {
		newErr := errors.New("insufficient balance")
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(newErr, http.StatusBadRequest))
		return
	}

	// check if campaign amount is greater than amount to be donated
	if float64(amount) > float64(campaign.Goal) {
		newErr := errors.New("amount to be donated is greater than campaign amount")
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(newErr, http.StatusBadRequest))
		return
	}

	// if goal is reached, close campaign
	if float64(campaign.TotalFunds) >= float64(campaign.Goal) {
		if err != nil {
			newErr := errors.New("campaign has closed")
			ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(newErr, http.StatusInternalServerError))
			return
		}
	}

	msg, err := defi.Donate(amount, idL, privateKey, address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, msg))
}

// @Summary Create campaign
// @Description Create campaign
// @Accept  json
// @Produce  json
// @Tags Campaigns
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param   title        formData   string    true  "Title"
// @Param   description        formData   string    true  "Description"
// @Param   goal        formData   string    true  "Goal"
// @Param   deadline        formData   string    true  "Deadline"
// @Param   category        formData   string    true  "Category"
// @Param   image        formData   file    true  "Image"
// @Success		200				{object}   string "hex"
// @Router /campaigns/create [post]
func (server *Server) createCampaign(ctx *gin.Context) {
	campaignImage, _, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	campaignTitle := ctx.Request.FormValue("title")
	campaignDescription := ctx.Request.FormValue("description")
	campaignGoal := ctx.Request.FormValue("goal")
	campaignDeadline := ctx.Request.FormValue("deadline")
	campaignCategory := ctx.Request.FormValue("category")

	// upload image to cloudinary

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

	uploadResult, err := utils.UploadAvatar(ctx, campaignImage, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	privateKey, address, err := crypt.DecryptPrivateKey(user.FilePath, server.config.PassPhase)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}
	// int64 to int
	// // add deadline to current time
	// dl := utils.ConvertToUnix(campaign.Deadline)

	// // convert string to float64
	goal, err := strconv.ParseFloat(campaignGoal, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	layoutString := "2006-01-02T15:04:05.000"
	// convert string to time
	deadline, err := time.Parse(layoutString, campaignDeadline)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	// check if deadline is less than current time
	if time.Now().After(deadline) {
		newErr := errors.New("deadline cannot be less than current time")
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(newErr, http.StatusBadRequest))
		return
	}

	campaigns, err := defi.CreateCampaign(campaignTitle, campaignCategory, campaignDescription, goal, deadline, uploadResult, privateKey, address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, campaigns))
}

// @Summary Get My Donations
// @Description Get My Donations
// @Accept  json
// @Produce  json
// @Tags Campaigns
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success		200				{object}    interfaces.DocSuccessResponse{data=[]crypt.Campaign}	"success"
// @Router /campaigns/donations [get]
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

	donations, err := defi.GetCampaignsByOwner(user.Address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, donations))
}

// @Summary Withdraw from campaign
// @Description Withdraw from campaign
// @Accept  json
// @Produce  json
// @Tags Campaigns
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param   data        body   interfaces.Withdraw[types.Post]    true  "Withdraw"
// @Success		200				{object}    string	"success"
// @Router /campaigns/withdraw [post]
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

	msg, err := defi.PayOut(withdraw.CampaignId, address, privateKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, msg))
}

// @Summary Get Current ETH Price
// @Description Get Current ETH Price
// @Accept  json
// @Produce  json
// @Tags Campaigns
// @Success		200				{object}    interfaces.DocSuccessResponse	"success"
// @Router /currentPrice [get]
func (server *Server) currentEthPrice(ctx *gin.Context) {
	price, err := defi.GetEthPrice()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, price))
}

// @Summary Get Campaign Categories
// @Description Get Campaign Categories
// @Accept  json
// @Produce  json
// @Tags Campaigns
// @Success		200				{object}    interfaces.DocSuccessResponse	"success"
// @Router /categories [get]
func (server *Server) getCategories(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload == nil {
		err := errors.New(interfaces.ErrUserNotFound)
		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
		return
	}

	_, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	campaigns, err := defi.GetCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	camps := make([]interfaces.CampaignCategory, len(campaigns))

	for i, campaign := range campaigns {
		camps[i] = interfaces.CampaignCategory{
			Name:        campaign.Name,
			Image:       campaign.Image,
			Description: campaign.Description,
			Id:          campaign.ID,
		}
	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, camps))

}

// @Summary Search Campaign by name
// @Description Search Campaign by name
// @Accept  json
// @Produce  json
// @Tags Campaigns
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param   data        body   interfaces.SearchCampaignRequest[types]    true  "SearchCampaignRequest"
// @Success		200				{object}    interfaces.DocSuccessResponse{data=[]interfaces.Campaigns}	"success"
// @Router /search [post]
func (server *Server) searchCampaignByName(ctx *gin.Context) {
	var req interfaces.SearchCampaignRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interfaces.ErrorResponse(err, http.StatusBadRequest))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload == nil {
		err := errors.New(interfaces.ErrUserNotFound)
		ctx.JSON(http.StatusNotFound, interfaces.ErrorResponse(err, http.StatusNotFound))
	}

	campaigns, err := defi.SearchCampaigns(req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
		return
	}

	camps := make([]interfaces.Campaigns, len(campaigns))

	for i, campaign := range campaigns {
		userInfo, _ := server.store.GetUserByAddress(ctx, campaign.Owner)

		totalNumber, err := defi.GetTotalDonationsByCampaignId(i)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
			return
		}

		//  Loop through the donators and amounts
		donators, amounts, _, err := defi.GetDonorsAddressesAndAmounts(i)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, interfaces.ErrorResponse(err, http.StatusInternalServerError))
			return
		}

		dons := make([]interfaces.DonorDetails, len(donators))

		// convert unix to time.Time
		deadline := time.Unix(int64(campaign.Deadline), 0)
		// if deadline is less than current time
		if deadline.Before(time.Now()) {
			// skip the current iteration
			continue
		} else {
			if len(donators) != 0 {
				for k, _ := range donators {
					getUser, _ := server.store.GetUserByAddress(ctx, donators[k])

					dons[k] = interfaces.DonorDetails{
						Amount:   (float64(amounts[k]) / 1e18),
						Donor:    donators[k],
						Image:    getUser.Avatar,
						Username: getUser.Username,
					}

					// wei to either
					goal := (float64(campaign.Goal) / 1e18)
					totalAmountDonated := (float64(campaign.TotalFunds) / 1e18)

					camps[i] = interfaces.Campaigns{
						CampaignType:       campaign.CampaignType,
						Title:              campaign.Title,
						Deadline:           time.Unix(int64(campaign.Deadline), 0),
						Description:        campaign.Description,
						Goal:               goal,
						Image:              campaign.Image,
						TotalAmountDonated: totalAmountDonated,
						TotalNumber:        totalNumber.Int64(),
						Owner:              campaign.Owner,
						ID:                 int(campaign.ID),
						Donations:          dons,
						User: []interfaces.UserResponseInfo{
							{
								Username: userInfo.Username,
								Email:    userInfo.Email,
								Address:  userInfo.Address,
								Avatar:   userInfo.Avatar,
							},
						},
					}
				}
			} else {
				goal := float64(campaign.Goal / 1e18)
				totalAmountDonated := float64(campaign.TotalFunds / 1e18)

				camps[i] = interfaces.Campaigns{
					CampaignType:       campaign.CampaignType,
					Title:              campaign.Title,
					Deadline:           time.Unix(int64(campaign.Deadline), 0),
					Description:        campaign.Description,
					Goal:               goal,
					Image:              campaign.Image,
					TotalAmountDonated: totalAmountDonated,
					TotalNumber:        totalNumber.Int64(),
					Owner:              campaign.Owner,
					ID:                 int(campaign.ID),
					Donations:          dons,
					User: []interfaces.UserResponseInfo{
						{
							Username: userInfo.Username,
							Email:    userInfo.Email,
							Address:  userInfo.Address,
							Avatar:   userInfo.Avatar,
						},
					},
				}
			}
		}

	}

	ctx.JSON(http.StatusOK, interfaces.Response(http.StatusOK, camps))
}
