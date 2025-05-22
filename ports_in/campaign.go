package ports_in

import (
	"ad-campaign-delivery/model"
	"context"
)

type CampaignService interface {
	Create(ctx context.Context, user model.Campaign) error
	Match(ctx context.Context, country model.Country, device model.Device, os model.OS) (*model.BidLookup, error)
}
