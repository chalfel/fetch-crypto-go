package usecase

import (
	"context"
	"fmt"

	"github.com/chalfel/fetch-crypto-go/src/crypto"
	"github.com/chalfel/fetch-crypto-go/src/external"
)

type FetchAllCryptoNameUsecase struct {
	CoinGecko        external.CoinGecko
	CryptoRepository crypto.CryptoRepository
}

func (fa *FetchAllCryptoNameUsecase) Fetch(ctx context.Context) error {
	geckoCryptos, err := fa.CoinGecko.ListCrypto()
	if err != nil {
		return err
	}
	fmt.Println(geckoCryptos)
	cryptos := []crypto.Crypto{}
	for _, geckoCrypto := range geckoCryptos {
		crypto := crypto.Crypto{Id: geckoCrypto.Id, Name: geckoCrypto.Name, Symbol: geckoCrypto.Symbol}
		cryptos = append(cryptos, crypto)
	}
	err = fa.CryptoRepository.CleanAndInsert(ctx, cryptos)
	return err
}
