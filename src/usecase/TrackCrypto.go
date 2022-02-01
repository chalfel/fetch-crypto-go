package usecase

import (
	"context"
	"fmt"

	"github.com/chalfel/fetch-crypto-go/src/crypto"
)

type TrackCryptoUsecase struct {
	CryptoMongoRepository crypto.CryptoMongoRepository
	CryptoTrackRepository crypto.CryptoTrackRepository
}

func (tc *TrackCryptoUsecase) Track(ctx context.Context) error {
	cryptos, err := tc.CryptoTrackRepository.GetAll(ctx)
	if err != nil {
		return err
	}
	for _, crypto := range cryptos {
		mongoCryptos, err := tc.CryptoMongoRepository.FindById(ctx, crypto.CryptoId)
		if err != nil {
			return err
		}
		for _, mongoCrypto := range mongoCryptos {
			fmt.Println(mongoCrypto)
		}
	}
	return nil
}
