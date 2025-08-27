package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Campaign represents the complete advertising campaign
// with targeting and budget information.
type Campaign struct {
	ID        string
	Country   Country
	Device    Device
	OS        OS
	Bid       decimal.Decimal
	Budget    decimal.Decimal
	Active    bool
	CreatedAt time.Time
	ExpiresAt time.Time
}

// Campaigns in the in memory implementation of campaigns storage.
// It is key-value store where each campaign ID maps to its corresponding Campaign data.
type Campaigns map[string]Campaign

// CampaignsLookup is a nested map structure that organizes campaigns by Country, Device, and OS,
// allowing efficient lookup and retrieval of bid data for targeted campaign delivery.
type CampaignsLookup map[Country]map[Device]map[OS][]BidLookup

// BidLookup contains the minimal campaign data required for
// bid delivery, including campaign ID and bid amount.
type BidLookup struct {
	ID  string
	Bid decimal.Decimal
}
