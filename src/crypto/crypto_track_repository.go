package crypto

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type CryptoTrackRepository struct {
	Db *pgx.Conn
}

func (cr *CryptoTrackRepository) Insert(ctx context.Context, crypto CryptoTrack) (*CryptoTrack, error) {
	err := cr.Db.QueryRow(ctx, "insert into crypto_track (crypto_id, user_email) values ($1, $2) returning id", crypto.CryptoId, crypto.UserEmail).Scan(&crypto.Id)
	if err != nil {
		return nil, err
	}

	return &crypto, nil
}

func (cr *CryptoTrackRepository) GetAll(ctx context.Context) ([]CryptoTrack, error) {
	cryptos := []CryptoTrack{}
	rows, err := cr.Db.Query(ctx, "SELECT * from crypto_track")
	for rows.Next() {
		crypto := CryptoTrack{}
		err := rows.Scan(&crypto.Id, &crypto.CryptoId, &crypto.UserEmail)
		if err != nil {
			return nil, err
		}
		cryptos = append(cryptos, crypto)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return cryptos, nil
}
