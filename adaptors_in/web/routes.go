package web

import (
	"net/http"

	"ad-campaign-delivery/ports_in"
)

func ConfigureCampaignRoutes(u ports_in.CampaignService, r *http.ServeMux) {
	campaignHandler := CampaignsHandler{UseCase: u}
	r.HandleFunc("POST /campaigns", campaignHandler.create)
	r.HandleFunc("POST /deliver", campaignHandler.match)
}
