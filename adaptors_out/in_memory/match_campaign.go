package in_memory

import (
	"ad-campaign-delivery/model"
	"ad-campaign-delivery/pkg"
	"context"
)

// MatchCampaign finds the highest available bid campaign to be delivered according to the
// informed params. Once the campaign is chosen, the bid value is deducted from the campaign budget.
func (r *CampaignRepository) MatchCampaign(ctx context.Context, country model.Country,
	device model.Device, os model.OS) (*model.BidLookup, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	orderedBids, ok := r.campaignsLookup[country][device][os]
	if !ok || len(orderedBids) == 0 {
		return nil, pkg.Errorf(pkg.ENOTFOUND, "no campaign found for %s, %s, %s", country, device, os)
	}

	for _, b := range orderedBids {
		if !r.campaigns[b.ID].Active {
			continue
		}
		r.deductBudget(b.ID)
		return &b, nil
	}
	// no campaign was found
	return nil, nil
}

func (r *CampaignRepository) deductBudget(campaignID string) {
	campaign, ok := r.campaigns[campaignID]
	if !ok {
		r.log.Error().
			Str("campaign_id", campaignID).
			Msg("attempt to deduct budget for non-existent campaign")
		return
	}

	campaign.Budget = campaign.Budget.Sub(campaign.Bid)

	if campaign.Bid.GreaterThan(campaign.Budget) {
		campaign.Active = false
	}

	r.campaigns[campaignID] = campaign
}
