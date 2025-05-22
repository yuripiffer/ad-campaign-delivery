package in_memory

import (
	"ad-campaign-delivery/model"
	"ad-campaign-delivery/pkg"
	"ad-campaign-delivery/pkg/logger"
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCampaignRepository_CreateCampaign(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(*CampaignRepository)
		campaign       model.Campaign
		lookupPosition int
		wantErr        error
	}{
		{
			name:  "create new campaign, creates new targeting keys and lookup",
			setup: func(r *CampaignRepository) {},
			campaign: model.Campaign{
				ID:        "camp1",
				Country:   model.France,
				Device:    model.Mobile,
				OS:        model.Android,
				Bid:       decimal.NewFromFloat(5),
				Budget:    decimal.NewFromFloat(100.5),
				Active:    true,
				CreatedAt: time.Now(),
			},
			lookupPosition: 0,
			wantErr:        nil,
		},
		{
			name: "should return error when campaign ID already exists",
			setup: func(r *CampaignRepository) {
				r.campaigns["camp1"] = model.Campaign{ID: "camp1"}
			},
			campaign: model.Campaign{
				ID:        "camp1",
				Country:   model.France,
				Device:    model.Mobile,
				OS:        model.Android,
				Bid:       decimal.NewFromFloat(5),
				Budget:    decimal.NewFromFloat(100.5),
				Active:    true,
				CreatedAt: time.Now(),
			},
			wantErr: pkg.Errorf(pkg.ECONFLICT, "campaign with ID camp1 already exists"),
		},
		{
			name: "create new campaign with highest bid, existing targeting keys",
			setup: func(r *CampaignRepository) {
				r.campaignsLookup = generateDefaultLookup()
			},
			campaign: model.Campaign{
				ID:        "camp1",
				Country:   model.France,
				Device:    model.Mobile,
				OS:        model.Android,
				Bid:       decimal.NewFromFloat(90.5),
				Budget:    decimal.NewFromFloat(1000.5),
				Active:    true,
				CreatedAt: time.Now(),
			},
			lookupPosition: 0,
			wantErr:        nil,
		},
		{
			name: "create new campaign, same value of bid already exists",
			setup: func(r *CampaignRepository) {
				r.campaignsLookup = generateDefaultLookup()
			},
			campaign: model.Campaign{
				ID:        "camp1",
				Country:   model.France,
				Device:    model.Mobile,
				OS:        model.Android,
				Bid:       decimal.NewFromFloat(30.1),
				Budget:    decimal.NewFromFloat(1000.5),
				Active:    true,
				CreatedAt: time.Now(),
			},
			lookupPosition: 3,
			wantErr:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logger.Init()
			repo := NewCampaignRepository(&l)
			tt.setup(repo)

			err := repo.CreateCampaign(context.Background(), tt.campaign)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)

			// Verify campaign was stored
			stored, exists := repo.campaigns[tt.campaign.ID]
			assert.True(t, exists)
			assert.Equal(t, tt.campaign, stored)

			// Verify lookup structure
			lookupSlice := repo.campaignsLookup[tt.campaign.Country][tt.campaign.Device][tt.campaign.OS]

			assert.NotEmpty(t, lookupSlice)

			foundInLookup := false
			for position, lookup := range lookupSlice {
				if lookup.ID == tt.campaign.ID {
					assert.Equal(t, tt.campaign.Bid, lookup.Bid)
					foundInLookup = true
					assert.Equal(t, tt.lookupPosition, position)
					break
				}
			}
			assert.True(t, foundInLookup)
		})
	}
}

func generateDefaultLookup() model.CampaignsLookup {
	return model.CampaignsLookup{
		model.France: {
			model.Mobile: {
				model.Android: {
					{ID: "a0", Bid: decimal.NewFromFloat(50)},
					{ID: "a1", Bid: decimal.NewFromFloat(30.1)},
					{ID: "a2", Bid: decimal.NewFromFloat(30.1)},
					{ID: "a3", Bid: decimal.NewFromFloat(20)},
				},
			},
		},
	}
}
