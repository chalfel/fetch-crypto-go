package crypto

type Crypto struct {
	Id     string
	Name   string
	Symbol string
}

type CryptoTrack struct {
	Id        string
	CryptoId  string
	UserEmail string
	CreatedAt string
	UpdatedAt string
}
