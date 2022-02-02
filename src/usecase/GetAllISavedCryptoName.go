package usecase

import (
	"context"

	"github.com/chalfel/fetch-crypto-go/src/crypto"
)

type GetAllCryptoNameUsecase struct {
	CryptoMongoRepository crypto.CryptoMongoRepository
}

func (fa *GetAllCryptoNameUsecase) GetAll(ctx context.Context, name string) ([]crypto.Crypto, error) {
	resp, err := fa.CryptoMongoRepository.SearchByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
