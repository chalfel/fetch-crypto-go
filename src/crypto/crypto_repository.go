package crypto

import (
	"context"
	"fmt"
	"strings"

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
	valueStrings := []string{}
	valueArgs := []interface{}{}

	for i, crypto := range cryptos {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))

		valueArgs = append(valueArgs, crypto.Id)
		valueArgs = append(valueArgs, crypto.Name)
		valueArgs = append(valueArgs, crypto.Symbol)
	}

	stmt := fmt.Sprintf("INSERT INTO crypto (id, name, symbol) VALUES %s", strings.Join(valueStrings, ","))
	_, err = tx.Exec(ctx, stmt, valueArgs...)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	return err
}
