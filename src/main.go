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
	if err != nil {
		log.Fatal("Error loading .env file")
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
		fmt.Println(cryptoTrack)
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
		fmt.Println(err)
		c.JSON(200, gin.H{"message": "Message fetched"})
	})

	r.GET("/sendEmail", func(c *gin.Context) {
		test := crypto.CryptoTrackRepository{Db: conn}
		resp, _ := test.GetAll(ctx)
		fmt.Println(resp)
		sendDto := external.SendDTO{CryptoId: "bitcoin", Email: "caiohalcsik@gmail.com", Percentage: 20.0, Value: 20.0}
		sendGrid := external.SendGrid{}
		sendGrid.Send(sendDto)
	})
	r.GET("/fetchCryptoName", func(c *gin.Context) {
		cryptoRepository := crypto.CryptoRepository{Db: conn}
		fetchAllCryptoNameUsecase := usecase.FetchAllCryptoNameUsecase{CoinGecko: external.CoinGecko{}, CryptoRepository: cryptoRepository}
		err := fetchAllCryptoNameUsecase.Fetch(ctx)
		fmt.Println(err)
	})
	r.GET("/crypto/all", func(c *gin.Context) {
		name, sent := c.GetQuery("name")
		if !sent {
			c.JSON(400, gin.H{"message": "name parameter is required"})
		}
		fmt.Println(name)
		cryptoRepository := crypto.CryptoRepository{Db: conn}
		getAllSavedCryptoName := usecase.GetAllCryptoNameUsecase{CryptoRepository: cryptoRepository}
		resp, err := getAllSavedCryptoName.GetAll(ctx, name)
		if err != nil {
			c.JSON(400, err)
			return
		}
		c.JSON(200, gin.H{"data": resp})
	})
	r.Run(":" + os.Getenv("PORT"))
}