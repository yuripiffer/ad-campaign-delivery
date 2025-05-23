package cron

import (
	"ad-campaign-delivery/ports_in"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type CampaignsHandler struct {
	UseCase ports_in.CampaignService
}

// CampaignExpirationChecker checks for expired campaigns every midnight and turn them to inactive.
func (h *CampaignsHandler) CampaignExpirationChecker(log zerolog.Logger) {
	cronjob := cron.New()

	// runs every day at 00:01
	_, err := cronjob.AddFunc("1 0 * * *", h.UseCase.DeactivateExpiredCampaigns)

	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to deactivate expired campaigns")
	}

	cronjob.Start()
}
