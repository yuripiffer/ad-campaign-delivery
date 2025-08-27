package in_memory

import (
	"context"
	"sync"
	"testing"

	"ad-campaign-delivery/model"
	"ad-campaign-delivery/pkg"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCampaignRepository_MatchCampaign(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *CampaignRepository
		country       model.Country
		device        model.Device
		os            model.OS
		wantBidLookup *model.BidLookup
		initialBudget decimal.Decimal
		wantErr       error
	}{
		{
			name: "campaign found, skips the first inactive higher bid, returns second bid, deduct budget",
			setup: func() *CampaignRepository {
				return &CampaignRepository{
					mu: sync.RWMutex{},
					campaigns: map[string]model.Campaign{
						"1": {ID: "1", Active: false},
						"2": {ID: "2", Active: true, Budget: decimal.NewFromFloat(1000),
							Bid: decimal.NewFromFloat(5)},
					},
					campaignsLookup: map[model.Country]map[model.Device]map[model.OS][]model.BidLookup{
						model.France: {
							model.Mobile: {
								model.Android: {
									{ID: "1", Bid: decimal.NewFromFloat(100)},
									{ID: "2", Bid: decimal.NewFromFloat(5)}},
							},
						},
					},
				}
			},
			country:       model.France,
			device:        model.Mobile,
			os:            model.Android,
			wantBidLookup: &model.BidLookup{ID: "2", Bid: decimal.NewFromFloat(5)},
			initialBudget: decimal.NewFromFloat(1000),
			wantErr:       nil,
		},
		{
			name: "no campaign found",
			setup: func() *CampaignRepository {
				return &CampaignRepository{
					mu:              sync.RWMutex{},
					campaigns:       map[string]model.Campaign{},
					campaignsLookup: map[model.Country]map[model.Device]map[model.OS][]model.BidLookup{},
				}
			},
			country:       model.France,
			device:        model.Mobile,
			os:            model.Android,
			wantBidLookup: nil,
			wantErr:       pkg.Errorf(pkg.ENOTFOUND, "no campaign found for FR, mobile, android"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setup()
			gotBidLookup, err := repo.MatchCampaign(context.Background(), tt.country, tt.device, tt.os)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantBidLookup, gotBidLookup)

			// make sure the budget was deducted
			if tt.wantBidLookup != nil {
				assert.Equal(t,
					tt.initialBudget.Sub(tt.wantBidLookup.Bid),
					repo.campaigns[tt.wantBidLookup.ID].Budget,
				)
			}
		})
	}
}
