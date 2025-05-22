package web

import (
	"ad-campaign-delivery/model"
	"ad-campaign-delivery/pkg"
	"ad-campaign-delivery/ports_in"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"net/http"
)

type CampaignsHandler struct {
	UseCase ports_in.CampaignService
}

type CampaignCreateRequest struct {
	ID      string          `json:"id"`
	Country string          `json:"country"`
	Device  string          `json:"device"`
	OS      string          `json:"os"`
	Bid     decimal.Decimal `json:"bid"`
	Budget  decimal.Decimal `json:"budget"`
}

// @Summary      Create a new campaign
// @Description  A campaign and a bid lookup will be created with the provided fields.
// @Tags         campaigns
// @Accept       json
// @Param        request  body  CampaignCreateRequest  true  "Campaign create request"
// @Success      201      "Campaign created (no content)"
// @Failure      400      {object}  pkg.ErrorResp
// @Failure      500      {object}  pkg.ErrorResp
// @Router       /campaigns [post]
func (h *CampaignsHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := CampaignCreateRequest{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		pkg.BadRequestResponse(w, r, fmt.Sprintf("invalid request payload: %v", err))
		return
	}

	if len(input.ID) == 0 {
		pkg.BadRequestResponse(w, r, "missing campaign ID")
		return
	}

	country, ok := model.Countries[input.Country]
	if !ok {
		pkg.BadRequestResponse(w, r, fmt.Sprintf("invalid country: %v", input.Country))
		return
	}

	device, ok := model.Devices[input.Device]
	if !ok {
		pkg.BadRequestResponse(w, r, fmt.Sprintf("invalid device: %v", input.Device))
		return
	}

	os, ok := model.OperationalSystems[input.OS]
	if !ok {
		pkg.BadRequestResponse(w, r, fmt.Sprintf("invalid os: %v", input.OS))
		return
	}

	if !input.Bid.IsPositive() {
		pkg.BadRequestResponse(w, r, fmt.Sprintf("invalid bid: %v", input.Bid))
		return
	}

	if input.Budget.IsNegative() {
		pkg.BadRequestResponse(w, r, fmt.Sprintf("invalid budget: %v", input.Budget))
		return
	}

	campaign := model.Campaign{
		ID:      input.ID,
		Country: country,
		Device:  device,
		OS:      os,
		Bid:     input.Bid,
		Budget:  input.Budget,
	}

	err = h.UseCase.Create(ctx, campaign)
	if err != nil {
		pkg.ErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type CampaignMatchRequest struct {
	Country string `json:"country"`
	Device  string `json:"device"`
	OS      string `json:"os"`
}

type CampaignMatchResponse struct {
	CampaignID string          `json:"campaign_id"`
	Bid        decimal.Decimal `json:"bid"`
}

// @Summary      Match a campaign
// @Description  Matches a campaign based on country, device, and OS, after validating consent.
// @Tags         campaigns
// @Accept       json
// @Produce      json
// @Param        X-Consent-String  header  string                  true  "Consent string"
// @Param        request           body    CampaignMatchRequest     true  "Campaign match request"
// @Success      200               {object} CampaignMatchResponse   "Matched campaign"
// @Success      204               "No matching campaign found"
// @Failure      400               {object} pkg.ErrorResp
// @Failure      500               {object} pkg.ErrorResp
// @Router       /campaigns/match [post]
func (h *CampaignsHandler) match(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	consentToken := r.Header.Get("X-Consent-String")
	if consentToken == "" {
		pkg.BadRequestResponse(w, r, "missing header X-Consent-String")
		return
	}
	consentVendorID := 1231 // let's assume Opti Digital Vendor ID is 1231 for now
	hasConsent, err := pkg.CheckConsent(consentToken, consentVendorID)
	if err != nil {
		pkg.ErrorResponse(w, r, err)
		return
	}
	if !hasConsent {
		pkg.BadRequestResponse(w, r, "invalid consent token")
		return
	}

	input := CampaignMatchRequest{}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		pkg.BadRequestResponse(w, r, fmt.Sprintf("invalid request payload: %v", err))
		return
	}

	country, ok := model.Countries[input.Country]
	if !ok {
		pkg.BadRequestResponse(w, r, fmt.Sprintf("invalid country: %v", input.Country))
		return
	}

	device, ok := model.Devices[input.Device]
	if !ok {
		pkg.BadRequestResponse(w, r, fmt.Sprintf("invalid device: %v", input.Device))
		return
	}

	os, ok := model.OperationalSystems[input.OS]
	if !ok {
		pkg.BadRequestResponse(w, r, fmt.Sprintf("invalid os: %v", input.OS))
		return
	}

	campaignMatch, err := h.UseCase.Match(ctx, country, device, os)
	if err != nil {
		pkg.ErrorResponse(w, r, err)
		return
	}

	if campaignMatch != nil {
		pkg.JsonResponse(w, r, http.StatusOK, CampaignMatchResponse{CampaignID: campaignMatch.ID, Bid: campaignMatch.Bid})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
