package in_memory

import (
	"ad-campaign-delivery/model"
	"ad-campaign-delivery/ports_out"
	"github.com/rs/zerolog"
	"sync"
)

type CampaignRepository struct {
	ports_out.CampaignRepository
	campaignsLookup model.CampaignsLookup
	campaigns       model.Campaigns
	mu              sync.RWMutex
	log             *zerolog.Logger
}

func NewCampaignRepository(log *zerolog.Logger) *CampaignRepository {
	return &CampaignRepository{
		campaignsLookup: model.CampaignsLookup{},
		campaigns:       model.Campaigns{},
		log:             log,
	}

}
