package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/chalfel/fetch-crypto-go/src/crypto"
	"github.com/chalfel/fetch-crypto-go/src/external"
	"github.com/chalfel/fetch-crypto-go/src/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if os.Getenv("ENVIRONMENT") == "development" {
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	cryptoMongoRepository := crypto.CryptoMongoRepository{Client: *client}
	if err != nil {
		fmt.Println(err)
	}

	r := gin.Default()
	r.POST("/crypto/track", func(c *gin.Context) {
		createCryptoTrack := usecase.CreateCryptoTrackUsecase{CryptoTrackRepository: crypto.CryptoTrackRepository{Db: conn}, CryptoMongoRepository: cryptoMongoRepository}
		var cryptoTrack crypto.CryptoTrack
		if err := c.BindJSON(&cryptoTrack); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		resp, err := createCryptoTrack.Create(ctx, cryptoTrack)
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": resp})
	})
	r.GET("/crypto", func(c *gin.Context) {
		fetchCryptoData := usecase.FetchCryptoDataUsecase{
			CoinGeckoClient:       external.CoinGecko{},
			CryptoMongoRepository: cryptoMongoRepository,
			CryptoTrackRepository: crypto.CryptoTrackRepository{Db: conn},
			SendGrid:              external.SendGrid{},
		}
		err := fetchCryptoData.Fetch(ctx)
		if err != nil {
			c.JSON(400, gin.H{"message": err})
			return
		}
		c.JSON(200, gin.H{"message": "Message fetched"})
	})

	r.GET("/crypto/name", func(c *gin.Context) {
		fetchAllCryptoNameUsecase := usecase.FetchAllCryptoNameUsecase{CoinGecko: external.CoinGecko{}, CryptoMongoRepository: cryptoMongoRepository}
		err := fetchAllCryptoNameUsecase.Fetch(ctx)
		if err != nil {
			c.JSON(400, gin.H{"message": err})
			return
		}
		c.JSON(200, gin.H{"message": "Crypto name is fetched"})
	})

	r.GET("/crypto/all", func(c *gin.Context) {
		name, sent := c.GetQuery("name")
		if !sent {
			c.JSON(400, gin.H{"message": "name parameter is required"})
		}
		fmt.Println(name)
		getAllSavedCryptoName := usecase.GetAllCryptoNameUsecase{CryptoMongoRepository: cryptoMongoRepository}
		resp, err := getAllSavedCryptoName.GetAll(ctx, name)
		if err != nil {
			c.JSON(400, err)
			return
		}
		c.JSON(200, gin.H{"data": resp})
	})

	r.GET("/crypto/send-email", func(c *gin.Context) {
		trackCryptoUsecase := usecase.TrackCryptoUsecase{
			CryptoMongoRepository: cryptoMongoRepository,
			CryptoTrackRepository: crypto.CryptoTrackRepository{Db: conn},
			SendGrid:              external.SendGrid{},
		}
		err := trackCryptoUsecase.Track(ctx)
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "E-mail sent"})
	})

	r.Run(":" + os.Getenv("PORT"))
}
