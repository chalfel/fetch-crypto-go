package usecase

import (
	"context"

	"github.com/chalfel/fetch-crypto-go/src/crypto"
	"github.com/chalfel/fetch-crypto-go/src/external"
)

type FetchAllCryptoNameUsecase struct {
	CoinGecko             external.CoinGecko
	CryptoMongoRepository crypto.CryptoMongoRepository
}

func (fa *FetchAllCryptoNameUsecase) Fetch(ctx context.Context) error {
	geckoCryptos, err := fa.CoinGecko.ListCrypto()
	if err != nil {
		return err
	}
	cryptos := []crypto.Crypto{}
	for _, geckoCrypto := range geckoCryptos {
		crypto := crypto.Crypto{Id: geckoCrypto.Id, Name: geckoCrypto.Name, Symbol: geckoCrypto.Symbol}
		cryptos = append(cryptos, crypto)
	}
	_, err = fa.CryptoMongoRepository.InsertAllCryptosName(ctx, cryptos)
	return err
}
