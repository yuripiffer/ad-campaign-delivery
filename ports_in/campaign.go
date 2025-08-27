package ports_in

import (
	"context"

	"ad-campaign-delivery/model"
)

//go:generate go run github.com/matryer/moq -out campaign_mock.go -stub . CampaignService
type CampaignService interface {
	Create(ctx context.Context, user model.Campaign, activeDays int) error
	Match(ctx context.Context, country model.Country, device model.Device, os model.OS) (*model.BidLookup, error)

	DeactivateExpiredCampaigns()
}
