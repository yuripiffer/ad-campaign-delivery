package app

import (
	"ad-campaign-delivery/adaptors_in/cron"
	"ad-campaign-delivery/adaptors_in/web"
	"ad-campaign-delivery/adaptors_out/in_memory"
	"ad-campaign-delivery/core/campaign"
	"ad-campaign-delivery/pkg/logger"
	"net/http"
	"time"
)

func Config() {
	time.Local = time.UTC
	log := logger.Init()

	campaignRepository := in_memory.NewCampaignRepository(&log)
	campaignService := campaign.NewService(campaignRepository)

	r := http.NewServeMux()
	web.ConfigureCampaignRoutes(campaignService, r)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}

	campaignCron := cron.CampaignsHandler{UseCase: campaignService}
	go campaignCron.CampaignExpirationChecker(log)
}
