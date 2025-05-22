package in_memory

func (r *CampaignRepository) deactivateCampaign(campaignID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	campaign, ok := r.campaigns[campaignID]
	if !ok {
		r.log.Error().
			Str("campaign_id", campaignID).
			Msg("attempt to deactivate non-existent campaign")
		return
	}

	campaign.Active = false
	r.campaigns[campaignID] = campaign
}
