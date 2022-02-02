package crypto

import (
	"context"
	"fmt"

	"github.com/chalfel/fetch-crypto-go/src/external"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CryptoMongoRepository struct {
	Client mongo.Client
}

func (cm *CryptoMongoRepository) InsertMany(ctx context.Context, cryptos []external.CoinGeckoCrypto) ([]external.CoinGeckoCrypto, error) {
	documents := make([]interface{}, len(cryptos))
	for i, crypto := range cryptos {
		documents[i] = crypto
	}
	_, err := cm.Client.Database("default").Collection("crypto").InsertMany(ctx, documents)
	if err != nil {
		return nil, err
	}
	return cryptos, nil
}

func (cm *CryptoMongoRepository) FindById(ctx context.Context, id string) ([]external.CoinGeckoCrypto, error) {
	filter := bson.M{
		"id": id,
	}
	response, err := cm.Client.Database("default").Collection("crypto").Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var decodedCryptoDocuments []external.CoinGeckoCrypto
	if err = response.All(context.TODO(), &decodedCryptoDocuments); err != nil {
		return nil, err
	}
	return decodedCryptoDocuments, nil
}

func (cm *CryptoMongoRepository) InsertAllCryptosName(ctx context.Context, cryptos []Crypto) ([]Crypto, error) {
	_, err := cm.Client.Database("default").Collection("crypto_name").DeleteMany(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	documents := make([]interface{}, len(cryptos))
	for i, crypto := range cryptos {
		documents[i] = crypto
	}
	_, err = cm.Client.Database("default").Collection("crypto_name").InsertMany(ctx, documents)
	if err != nil {
		return nil, err
	}

	return cryptos, nil
}

func (cm *CryptoMongoRepository) SearchByName(ctx context.Context, name string) ([]Crypto, error) {
	regexPattern := fmt.Sprintf(".*%s.*", name)
	filter := bson.D{primitive.E{Key: "name", Value: primitive.Regex{Pattern: regexPattern, Options: ""}}}

	response, err := cm.Client.Database("default").Collection("crypto_name").Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var decodedCryptoDocuments []Crypto
	if err = response.All(ctx, &decodedCryptoDocuments); err != nil {
		return nil, err
	}
	return decodedCryptoDocuments, nil
}
