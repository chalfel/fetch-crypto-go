package crypto

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type cryptoRepository interface {
	GetAll(ctx context.Context, name string) ([]Crypto, error)
	CleanAndInsert(ctx context.Context, cryptos []Crypto) error
}

type CryptoRepository struct {
	Db *pgx.Conn
}

func (cr *CryptoRepository) GetAll(ctx context.Context, name string) ([]Crypto, error) {
	rows, err := cr.Db.Query(ctx, "select * from crypto where name like '%' || $1 || '%'", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cryptos := []Crypto{}
	for rows.Next() {
		crypto := Crypto{}
		rows.Scan(&crypto.Id, &crypto.Name, &crypto.Symbol)
		cryptos = append(cryptos, crypto)
	}
	return cryptos, nil
}

func (cr *CryptoRepository) CleanAndInsert(ctx context.Context, cryptos []Crypto) error {
	tx, err := cr.Db.Begin(ctx)
	if err != nil {
		return err
	}
	for _, crypto := range cryptos {
		tx.QueryRow(ctx, "insert into crypto (id, name, symbol) values ($1, $2, $3)", &crypto.Id, &crypto.Name, &crypto.Symbol).Scan()
	}
	err = tx.Commit(ctx)
	return err
}
