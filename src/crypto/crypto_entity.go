package crypto

import "time"

type Crypto struct {
	Id     string
	Name   string
	Symbol string
}

type CryptoTrack struct {
	Id        string
	CryptoId  string
	UserEmail string
	CreatedAt time.Time
	UpdatedAt time.Time
}
