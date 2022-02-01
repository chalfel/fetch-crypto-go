package usecase

import (
	"context"

	"github.com/chalfel/fetch-crypto-go/src/crypto"
)

type GetAllCryptoNameUsecase struct {
	CryptoRepository crypto.CryptoRepository
}

func (fa *GetAllCryptoNameUsecase) GetAll(ctx context.Context, name string) ([]crypto.Crypto, error) {
	resp, err := fa.CryptoRepository.GetAll(ctx, name)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
