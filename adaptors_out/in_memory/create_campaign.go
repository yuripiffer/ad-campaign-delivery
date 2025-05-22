package in_memory

import (
	"ad-campaign-delivery/model"
	"ad-campaign-delivery/pkg"
	"context"
)

// CreateCampaign inserts a new campaign into the in-memory store.
// All data must have already been validated in the domain.
func (r *CampaignRepository) CreateCampaign(ctx context.Context, campaign model.Campaign) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.campaigns[campaign.ID]; ok {
		return pkg.Errorf(pkg.ECONFLICT, "campaign with ID %s already exists", campaign.ID)
	}
	r.campaigns[campaign.ID] = campaign

	if _, ok := r.campaignsLookup[campaign.Country][campaign.Device][campaign.OS]; !ok {
		r.createTargetingKeys(campaign)
	}
	r.insertBidInLookup(campaign)
	return nil
}

// createTargetingKeys initializes the targeting keys in the campaignsLookup map
func (r *CampaignRepository) createTargetingKeys(campaign model.Campaign) {
	if _, ok := r.campaignsLookup[campaign.Country]; !ok {
		r.campaignsLookup[campaign.Country] = make(map[model.Device]map[model.OS][]model.BidLookup)
	}
	if _, ok := r.campaignsLookup[campaign.Country][campaign.Device]; !ok {
		r.campaignsLookup[campaign.Country][campaign.Device] = make(map[model.OS][]model.BidLookup)
	}
	if _, ok := r.campaignsLookup[campaign.Country][campaign.Device][campaign.OS]; !ok {
		r.campaignsLookup[campaign.Country][campaign.Device][campaign.OS] = []model.BidLookup{}
	}
}

// insertBidInLookup guarantees that older campaigns with the same bid should be selected
// first due to how to slice is populated
func (r *CampaignRepository) insertBidInLookup(campaign model.Campaign) {
	orderedBids := r.campaignsLookup[campaign.Country][campaign.Device][campaign.OS]
	newBid := model.BidLookup{
		ID:  campaign.ID,
		Bid: campaign.Bid,
	}

	// Binary search for insertion index (descending order)
	low, high := 0, len(orderedBids)
	for low < high {
		mid := (low + high) / 2
		if orderedBids[mid].Bid.GreaterThanOrEqual(newBid.Bid) {
			low = mid + 1
		} else {
			high = mid
		}
	}

	// Insert new bidLookup at the determined index
	r.campaignsLookup[campaign.Country][campaign.Device][campaign.OS] = append(
		orderedBids[:low],
		append([]model.BidLookup{newBid}, orderedBids[low:]...)...,
	)
}
