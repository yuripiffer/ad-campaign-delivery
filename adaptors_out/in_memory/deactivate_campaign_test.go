package in_memory

import (
	"testing"
	"time"

	"ad-campaign-delivery/model"
)

func TestDeactivateExpiredCampaigns(t *testing.T) {
	now := time.Now()
	campaigns := model.Campaigns{
		"1": {
			ID:        "1",
			Active:    true,
			ExpiresAt: now.Add(-time.Hour),
		},
		"2": {
			ID:        "2",
			Active:    true,
			ExpiresAt: now.Add(time.Hour),
		},
		"3": {
			ID:        "3",
			Active:    false,
			ExpiresAt: now.Add(-2 * time.Hour),
		},
	}

	repo := &CampaignRepository{
		campaigns: campaigns,
	}

	repo.DeactivateExpiredCampaigns()

	if campaigns["1"].Active {
		t.Errorf("Expected campaign 1 to be deactivated")
	}
	if !campaigns["2"].Active {
		t.Errorf("Expected campaign 2 to remain active")
	}
	if campaigns["3"].Active {
		t.Errorf("Expected campaign 3 to remain inactive")
	}
}
