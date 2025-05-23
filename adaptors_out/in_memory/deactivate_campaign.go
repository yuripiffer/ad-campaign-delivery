package in_memory

import (
	"time"
)

func (r *CampaignRepository) DeactivateExpiredCampaigns() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, c := range r.campaigns {
		if c.Active && c.ExpiresAt.Before(time.Now()) {
			c.Active = false
			r.campaigns[id] = c
		}
	}
}
