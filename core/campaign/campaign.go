package campaign

import (
	"ad-campaign-delivery/model"
	"ad-campaign-delivery/ports_in"
	"ad-campaign-delivery/ports_out"
	"context"
	"time"
)

type Service struct {
	ports_in.CampaignService
	campaignRepository ports_out.CampaignRepository
}

func NewService(campaignRepository ports_out.CampaignRepository) *Service {
	return &Service{
		campaignRepository: campaignRepository,
	}
}

func (s *Service) Create(ctx context.Context, campaign model.Campaign) error {
	campaign.CreatedAt = time.Now()
	if campaign.Budget.GreaterThanOrEqual(campaign.Bid) {
		campaign.Active = true
	}

	return s.campaignRepository.CreateCampaign(ctx, campaign)
}

func (s *Service) Match(ctx context.Context, country model.Country, device model.Device, os model.OS) (*model.BidLookup, error) {
	return s.campaignRepository.MatchCampaign(ctx, country, device, os)
}
