package web

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ad-campaign-delivery/model"
	"ad-campaign-delivery/pkg"
	"ad-campaign-delivery/ports_in"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCampaignsHandler_Create(t *testing.T) {
	tests := []struct {
		name         string
		input        CampaignCreateRequest
		callCreate   bool
		createErr    error
		expectedCode int
		expectedBody string
	}{
		{
			name: "successful creation",
			input: CampaignCreateRequest{
				ID:         "camp123",
				Country:    "FR",
				Device:     "mobile",
				OS:         "android",
				Bid:        decimal.NewFromFloat(1.5),
				Budget:     decimal.NewFromFloat(100),
				ActiveDays: 30,
			},
			callCreate:   true,
			createErr:    nil,
			expectedCode: http.StatusCreated,
		},
		{
			name: "missing ID",
			input: CampaignCreateRequest{
				Country: "FR",
				Device:  "mobile",
				OS:      "android",
				Bid:     decimal.NewFromFloat(1.5),
				Budget:  decimal.NewFromFloat(100),
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "missing campaign ID",
		},
		{
			name: "invalid country",
			input: CampaignCreateRequest{
				ID:      "camp123",
				Country: "invalid_country",
				Device:  "mobile",
				OS:      "android",
				Bid:     decimal.NewFromFloat(1.5),
				Budget:  decimal.NewFromFloat(100),
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid country: invalid_country",
		},
		{
			name: "invalid device",
			input: CampaignCreateRequest{
				ID:      "camp123",
				Country: "FR",
				Device:  "invalid_device",
				OS:      "android",
				Bid:     decimal.NewFromFloat(1.5),
				Budget:  decimal.NewFromFloat(100),
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid device: invalid_device",
		},
		{
			name: "invalid os",
			input: CampaignCreateRequest{
				ID:      "camp123",
				Country: "FR",
				Device:  "mobile",
				OS:      "invalid_os",
				Bid:     decimal.NewFromFloat(1.5),
				Budget:  decimal.NewFromFloat(100),
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid os: invalid_os",
		},
		{
			name: "invalid bid with negative value",
			input: CampaignCreateRequest{
				ID:      "camp123",
				Country: "FR",
				Device:  "mobile",
				OS:      "android",
				Bid:     decimal.NewFromFloat(-1.5),
				Budget:  decimal.NewFromFloat(100),
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid bid: -1.5",
		},
		{
			name: "invalid bid with value 0",
			input: CampaignCreateRequest{
				ID:      "camp123",
				Country: "FR",
				Device:  "mobile",
				OS:      "android",
				Bid:     decimal.NewFromFloat(0),
				Budget:  decimal.NewFromFloat(100),
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid bid: 0",
		},
		{
			name: "invalid budget",
			input: CampaignCreateRequest{
				ID:      "camp123",
				Country: "FR",
				Device:  "mobile",
				OS:      "android",
				Bid:     decimal.NewFromFloat(1.5),
				Budget:  decimal.NewFromFloat(-100),
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid budget: -100",
		},
		{
			name: "error from Create method in domain",
			input: CampaignCreateRequest{
				ID:      "camp123",
				Country: "FR",
				Device:  "mobile",
				OS:      "android",
				Bid:     decimal.NewFromFloat(1.5),
				Budget:  decimal.NewFromFloat(100),
			},
			callCreate:   true,
			createErr:    pkg.Errorf(pkg.ECONFLICT, "campaign with ID camp123 already exists"),
			expectedCode: http.StatusConflict,
			expectedBody: "campaign with ID camp123 already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaignServiceMock := &ports_in.CampaignServiceMock{
				CreateFunc: func(ctx context.Context, campaign model.Campaign, activeDays int) error {
					assert.Equal(t, tt.input.ID, campaign.ID)
					assert.Equal(t, model.Countries[tt.input.Country], campaign.Country)
					assert.Equal(t, model.Devices[tt.input.Device], campaign.Device)
					assert.Equal(t, model.OperationalSystems[tt.input.OS], campaign.OS)
					assert.True(t, tt.input.Bid.Equal(campaign.Bid))
					assert.True(t, tt.input.Budget.Equal(campaign.Budget))
					assert.Equal(t, tt.input.ActiveDays, activeDays)
					return tt.createErr
				},
			}
			handler := CampaignsHandler{UseCase: campaignServiceMock}

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/campaigns", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			handler.create(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestCampaignsHandler_Match(t *testing.T) {

	// String for TCF v2 format with valid consents
	validConsentString := "CQMGLkAQMGLkABcAKEFRBbFgAP_gAEPgAAqIJnkR_C9MQWFjcT51AfskaYxHxgACo" +
		"EQgBACJgygBCAPA8IQEwGAYIAxAAqAKAAAAoiRBAAAlCAhQAAAAQAAAACCMAEAAAAAAIKBAgAARAgEACAhB" +
		"GQAAEAAAAIBBABAAgAAEQBoAQBAAAAAAAAAgAAAgAACBAAAIAAAAAAEAAAAIAEgAAAAAAAAAAAAAAlAIAAA" +
		"IAAAAAAAAAAAIJngAmChEQAFgQAhAAGEECABQRgAAAAAgAACBggAACAAA4AQAUGAAAAAAAAAIAAAAggABAAA" +
		"BAAhAAAAAQAAAAAAIAAAAAAAAACBAAAABAAAAAAgAAQAAAAAAAABAABAAgAAAABAAQBAAAAAgAAAAAAAAAAC" +
		"AAAAAAAAAAAEAAAAIAEAAAAAAAAAAAAAAAAAIAAAAAAAAAAAAAAAAAAA"

	//Valid TCF v2 format, but missing consent
	missingConsentString := "COtybn4Otybn4AcABBENAPCIAEBAAECAAIAAAAAAAAAAAgAA.YAAAAAAAAAAA"

	successfulMatch := `{
	"campaign_id": "camp123",
	"bid": "1.5"
}
`
	tests := []struct {
		name              string
		consentToken      string
		input             CampaignMatchRequest
		callMatch         bool
		mockMatchResponse *model.BidLookup
		mockMatchError    error
		expectedCode      int
		expectedBody      string
	}{
		{
			name:         "successful match",
			consentToken: validConsentString,
			input: CampaignMatchRequest{
				Country: "FR",
				Device:  "mobile",
				OS:      "android",
			},
			callMatch: true,
			mockMatchResponse: &model.BidLookup{
				ID:  "camp123",
				Bid: decimal.NewFromFloat(1.5),
			},
			expectedCode: http.StatusOK,
			expectedBody: successfulMatch,
		},
		{
			name:         "success, no match found",
			consentToken: validConsentString,
			input: CampaignMatchRequest{
				Country: "FR",
				Device:  "mobile",
				OS:      "android",
			},
			callMatch:         true,
			mockMatchResponse: nil,
			expectedCode:      http.StatusNoContent,
		},
		{
			name:         "missing consent token",
			consentToken: "",
			input: CampaignMatchRequest{
				Country: "FR",
				Device:  "mobile",
				OS:      "android",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "missing header X-Consent-String",
		},
		{
			name:         "failed to parse token via iabconsent V2",
			consentToken: "thisIsNotAValidV2ConsentString",
			input: CampaignMatchRequest{
				Country: "FR",
				Device:  "mobile",
				OS:      "android",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "failed to parse consent string",
		},
		{
			name:         "missing consent",
			consentToken: missingConsentString,
			input: CampaignMatchRequest{
				Country: "FR",
				Device:  "mobile",
				OS:      "android",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid consent token",
		},
		{
			name:         "invalid country",
			consentToken: validConsentString,
			input: CampaignMatchRequest{
				Country: "invalid_country",
				Device:  "mobile",
				OS:      "android",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid country: invalid_country",
		},
		{
			name:         "invalid device",
			consentToken: validConsentString,
			input: CampaignMatchRequest{
				Country: "FR",
				Device:  "invalid_device",
				OS:      "android",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid device: invalid_device",
		},
		{
			name:         "invalid os",
			consentToken: validConsentString,
			input: CampaignMatchRequest{
				Country: "FR",
				Device:  "mobile",
				OS:      "invalid_os",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid os: invalid_os",
		},
		{
			name:         "error from Match method in domain",
			consentToken: validConsentString,
			input: CampaignMatchRequest{
				Country: "FR",
				Device:  "mobile",
				OS:      "android",
			},
			callMatch:      true,
			mockMatchError: pkg.Errorf(pkg.EINTERNAL, "internal error"),
			expectedCode:   http.StatusInternalServerError,
			expectedBody:   "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaignServiceMock := &ports_in.CampaignServiceMock{
				MatchFunc: func(ctx context.Context, country model.Country, device model.Device,
					os model.OS) (*model.BidLookup, error) {
					return tt.mockMatchResponse, tt.mockMatchError
				},
			}
			handler := CampaignsHandler{UseCase: campaignServiceMock}

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/deliver", bytes.NewBuffer(body))
			if tt.consentToken != "" {
				req.Header.Set("X-Consent-String", tt.consentToken)
			}
			rec := httptest.NewRecorder()

			handler.match(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}
