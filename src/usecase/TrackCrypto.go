package usecase

import (
	"context"

	"github.com/chalfel/fetch-crypto-go/src/crypto"
	"github.com/chalfel/fetch-crypto-go/src/external"
)

type TrackCryptoUsecase struct {
	CryptoMongoRepository crypto.CryptoMongoRepository
	CryptoTrackRepository crypto.CryptoTrackRepository
	SendGrid              external.SendGrid
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
		var percentage, value float64
		percentage = mongoCryptos[0].PriceChangePercentage1hInCurrency
		value = mongoCryptos[0].CurrentPrice
		sendDTO := external.SendDTO{
			CryptoId:   crypto.CryptoId,
			Email:      crypto.UserEmail,
			Percentage: percentage,
			Value:      value,
		}
		err = tc.SendGrid.Send(sendDTO)
		if err != nil {
			return err
		}
	}
	return err
}
