package external

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

type CoinGecko struct{}
type ListCryptoObject struct {
	Id     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}
type Roi struct {
	Times      float64 `json:"times"`
	Currency   string  `json:"currency"`
	Percentage float64 `json:"percentage"`
}

type CoinGeckoCrypto struct {
	Id                                string  `json:"id"`
	Symbol                            string  `json:"symbol"`
	Name                              string  `json:"name"`
	Image                             string  `json:"image"`
	CurrentPrice                      float64 `json:"current_price"`
	MarketCap                         float64 `json:"market_cap"`
	MarketCapRank                     float64 `json:"market_cap_rank"`
	FullyDilutedValuation             float64 `json:"fully_diluted_valuation"`
	TotalVolume                       float64 `json:"total_volume"`
	High24h                           float64 `json:"high_24h"`
	Low24h                            float64 `json:"low_24h"`
	PriceChange24h                    float64 `json:"price_change_24h"`
	PriceChangePercentage24h          float64 `json:"price_change_percentage_24h"`
	MarketCapChange24h                float64 `json:"market_cap_change_24h"`
	MarketCapChangePercentage24h      float64 `json:"market_cap_change_percentage_24h"`
	CirculatingSupply                 float64 `json:"circulating_supply"`
	TotalSupply                       float64 `json:"total_supply"`
	MaxSupply                         float64 `json:"max_supply"`
	Ath                               float64 `json:"ath"`
	AthChangePercentage               float64 `json:"ath_change_percentage"`
	AthDate                           string  `json:"ath_date"`
	Atl                               float64 `json:"atl"`
	AtlChangePercentage               float64 `json:"atl_change_percentage"`
	PriceChangePercentage1hInCurrency float64 `json:"price_change_percentage_1h_in_currency"`
	AtlDate                           string  `json:"atl_date"`
	Roi                               Roi
	LastUpdated                       string `json:"last_updated"`
}

func (c *CoinGecko) ListCrypto() ([]ListCryptoObject, error) {
	response, err := httpClient.Get("https://api.coingecko.com/api/v3/coins/list")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	cryptoListArray := []ListCryptoObject{}
	err = json.NewDecoder(response.Body).Decode(&cryptoListArray)
	if err != nil {
		return nil, err
	}
	return cryptoListArray, nil
}
func (c *CoinGecko) GetCoinMarket(page int, allCryptoMarketArray *[]CoinGeckoCrypto) {

	response, err := httpClient.Get("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&price_change_percentage=1h&per_page=250&page=" + fmt.Sprint(page))
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()
	cryptoMarketArray := []CoinGeckoCrypto{}
	err = json.NewDecoder(response.Body).Decode(&cryptoMarketArray)
	if err != nil {
		fmt.Println(err)
	}
	*allCryptoMarketArray = append(*allCryptoMarketArray, cryptoMarketArray...)

}
func (c *CoinGecko) GetAllMarketInfo() []CoinGeckoCrypto {
	cryptoArr, err := c.ListCrypto()
	if err != nil {
		log.Println(err)
	}
	cryptoPages := math.Ceil(float64(len(cryptoArr)) / float64(250))

	allCryptoMarketArray := make([]CoinGeckoCrypto, len(cryptoArr))
	for i := 0; i < int(cryptoPages); i++ {
		go c.GetCoinMarket(i, &allCryptoMarketArray)
	}

	// page := 1
	// stop := false
	// allCryptoMarketArray := []CoinGeckoCrypto{}
	// for !stop {
	// 	response, err := httpClient.Get("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&price_change_percentage=1h&per_page=250&page=" + fmt.Sprint(page))
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	defer response.Body.Close()
	// 	cryptoMarketArray := []CoinGeckoCrypto{}
	// 	err = json.NewDecoder(response.Body).Decode(&cryptoMarketArray)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	if len(cryptoMarketArray) == 0 {
	// 		stop = true
	// 	} else {
	// 		allCryptoMarketArray = append(allCryptoMarketArray, cryptoMarketArray...)
	// 		page = page + 1
	// 	}
	// }
	return allCryptoMarketArray
}
