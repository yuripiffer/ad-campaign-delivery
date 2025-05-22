package model

import (
	"github.com/shopspring/decimal"
	"time"
)

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

type Campaigns map[string]Campaign

type CampaignsLookup map[Country]map[Device]map[OS][]BidLookup

type BidLookup struct {
	ID  string
	Bid decimal.Decimal
}
