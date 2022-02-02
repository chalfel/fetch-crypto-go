package usecase

import (
	"context"

	"github.com/chalfel/fetch-crypto-go/src/crypto"
	"github.com/chalfel/fetch-crypto-go/src/external"
)

type FetchCryptoDataUsecase struct {
	CoinGeckoClient       external.CoinGecko
	CryptoMongoRepository crypto.CryptoMongoRepository
	CryptoTrackRepository crypto.CryptoTrackRepository
	SendGrid              external.SendGrid
}

func (fc *FetchCryptoDataUsecase) Fetch(ctx context.Context) error {
	cryptos := fc.CoinGeckoClient.GetAllMarketInfo()
	_, err := fc.CryptoMongoRepository.InsertMany(ctx, cryptos)

	if err != nil {
		return err
	}

	cryptoTrackings, err := fc.CryptoTrackRepository.GetAll(ctx)
	if err != nil {
		return err
	}
	for _, track := range cryptoTrackings {
		var percentage, value float64
		for _, crypto := range cryptos {
			if crypto.Id == track.CryptoId {
				percentage = crypto.PriceChangePercentage1hInCurrency
				value = crypto.CurrentPrice
			}
		}
		sendDTO := external.SendDTO{
			CryptoId:   track.CryptoId,
			Email:      track.UserEmail,
			Percentage: percentage,
			Value:      value,
		}
		fc.SendGrid.Send(sendDTO)
	}
	return err

}
