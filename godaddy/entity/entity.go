package entity

import "time"

type Secret struct {
	ApiKey    string
	ApiSecret string
}

type Domain struct {
	CreatedAt              time.Time `json:"createdAt"`
	Domain                 string    `json:"domain"`
	DomainID               int       `json:"domainId"`
	ExpirationProtected    bool      `json:"expirationProtected"`
	Expires                time.Time `json:"expires"`
	ExposeWhois            bool      `json:"exposeWhois"`
	HoldRegistrar          bool      `json:"holdRegistrar"`
	Locked                 bool      `json:"locked"`
	NameServers            any       `json:"nameServers"`
	Privacy                bool      `json:"privacy"`
	RegistrarCreatedAt     time.Time `json:"registrarCreatedAt"`
	RenewAuto              bool      `json:"renewAuto"`
	RenewDeadline          time.Time `json:"renewDeadline"`
	Renewable              bool      `json:"renewable"`
	Status                 string    `json:"status"`
	TransferAwayEligibleAt time.Time `json:"transferAwayEligibleAt"`
	TransferProtected      bool      `json:"transferProtected"`
}

type Record struct {
	Data string `json:"data"`
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
	Type string `json:"type"`
}
