package ports_out

import (
	"ad-campaign-delivery/model"
	"context"
)
//go:generate go run github.com/matryer/moq -out campaign_mock.go -stub . CampaignRepository
type CampaignRepository interface {
	CreateCampaign(ctx context.Context, campaign model.Campaign) error
	MatchCampaign(ctx context.Context, country model.Country, device model.Device, os model.OS) (*model.BidLookup, error)
	DeactivateExpiredCampaigns()
}
