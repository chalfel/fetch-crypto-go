package usecase

import (
	"context"
	"errors"

	"github.com/chalfel/fetch-crypto-go/src/crypto"
)

type CreateCryptoTrackUsecase struct {
	CryptoTrackRepository crypto.CryptoTrackRepository
	CryptoMongoRepository crypto.CryptoMongoRepository
}

func (cc *CreateCryptoTrackUsecase) Create(ctx context.Context, cryptoTrack crypto.CryptoTrack) (*crypto.CryptoTrack, error) {
	response, err := cc.CryptoMongoRepository.FindById(ctx, cryptoTrack.CryptoId)

	if err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return nil, errors.New("crypto not found")
	}

	insertedCryptoTrack, err := cc.CryptoTrackRepository.Insert(ctx, cryptoTrack)

	if err != nil {
		return nil, err
	}

	return insertedCryptoTrack, nil
}
