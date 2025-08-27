package campaign

import (
	"context"
	"testing"
	"time"

	"ad-campaign-delivery/model"
	"ad-campaign-delivery/ports_out"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCampaignService_Create(t *testing.T) {
	timeNowMock := time.Now()
	expiresAtMock := timeNowMock.AddDate(0, 0, 30)
	tests := []struct {
		name              string
		inputCampaign     model.Campaign
		activeDays        int
		CampaignToPersist model.Campaign
		wantErr           bool
	}{
		{
			name: "budget is lower than bid",
			inputCampaign: model.Campaign{
				ID: "1", Country: model.France, Device: model.Mobile, OS: model.Android,
				Bid: decimal.NewFromFloat(50), Budget: decimal.NewFromFloat(10)},

			CampaignToPersist: model.Campaign{
				ID: "1", Country: model.France, Device: model.Mobile, OS: model.Android,
				Bid: decimal.NewFromFloat(50), Budget: decimal.NewFromFloat(10),
				Active: false, CreatedAt: timeNowMock},
		},

		{
			name: "budget is equal to bid, with expiration day",
			inputCampaign: model.Campaign{
				ID: "1", Country: model.France, Device: model.Mobile, OS: model.Android,
				Bid: decimal.NewFromFloat(10), Budget: decimal.NewFromFloat(10)},

			activeDays: 30,

			CampaignToPersist: model.Campaign{
				ID: "1", Country: model.France, Device: model.Mobile, OS: model.Android,
				Bid: decimal.NewFromFloat(10), Budget: decimal.NewFromFloat(10),
				Active: true, CreatedAt: timeNowMock, ExpiresAt: expiresAtMock},
		},
		{
			name: "budget is higher than bid value",
			inputCampaign: model.Campaign{
				ID: "1", Country: model.France, Device: model.Mobile, OS: model.Android,
				Bid: decimal.NewFromFloat(10), Budget: decimal.NewFromFloat(1000)},

			CampaignToPersist: model.Campaign{
				ID: "1", Country: model.France, Device: model.Mobile, OS: model.Android,
				Bid: decimal.NewFromFloat(10), Budget: decimal.NewFromFloat(1000),
				Active: true, CreatedAt: timeNowMock},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaignRepo := &ports_out.CampaignRepositoryMock{
				CreateCampaignFunc: func(ctx context.Context, c model.Campaign) error {

					// add a bit of time tolerance to avoid time differences
					if c.CreatedAt.Before(timeNowMock.Add(1*time.Minute)) &&
						c.CreatedAt.After(timeNowMock.Add(-1*time.Minute)) {
						c.CreatedAt = timeNowMock
					}
					if c.ExpiresAt.Before(expiresAtMock.Add(1*time.Minute)) &&
						c.ExpiresAt.After(expiresAtMock.Add(-1*time.Minute)) {
						c.ExpiresAt = expiresAtMock
					}

					assert.Equal(t, tt.CampaignToPersist, c)
					return nil
				},
			}

			service := NewService(campaignRepo)
			err := service.Create(context.Background(), tt.inputCampaign, tt.activeDays)
			assert.NoError(t, err)
		})
	}
}

func TestCampaignService_Match(t *testing.T) {
	tests := []struct {
		name      string
		country   model.Country
		Device    model.Device
		OS        model.OS
		bidLookup *model.BidLookup
	}{
		{
			name:    "delivers bid",
			country: model.France,
			Device:  model.Mobile,
			OS:      model.Android,
			bidLookup: &model.BidLookup{
				ID:  "123",
				Bid: decimal.NewFromFloat(5),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			campaignRepo := &ports_out.CampaignRepositoryMock{
				MatchCampaignFunc: func(ctx context.Context, country model.Country,
					device model.Device, os model.OS) (*model.BidLookup, error) {

					assert.Equal(t, tt.country, country)
					assert.Equal(t, tt.Device, device)
					assert.Equal(t, tt.OS, os)

					return tt.bidLookup, nil
				},
			}

			service := NewService(campaignRepo)
			bidLookup, err := service.Match(context.Background(), tt.country, tt.Device, tt.OS)
			assert.NoError(t, err)
			if tt.bidLookup != nil {
				assert.Equal(t, tt.bidLookup, bidLookup)
			}
		})
	}
}
