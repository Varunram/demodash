package main

import (
	"encoding/json"
	"log"
	"math"
	"sync"

	errors "github.com/pkg/errors"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
)

// package tickers implements handlers for getting price from cryptocurrencyu markets

// we take the three largest (no wash trading) markets for BTC USD and return their weighted average
// to arrive at the price for BTC-USD. This price is indicative and not final since there will be latency
// involved between price display and trade finality.

// BinanceReqBTC is the binance ticker from the API
var BinanceReqBTC = "https://api.binance.com/api/v1/ticker/price?symbol=BTCUSDT"

var BinanceReqETH = "https://api.binance.com/api/v1/ticker/price?symbol=ETHUSDT"

// CoinbaseReqBTC is the coinbase ticker from the API
var CoinbaseReqBTC = "https://api.pro.coinbase.com/products/BTC-USD/ticker"

// CoinbaseReqETH is the coinbase ticker from the API
var CoinbaseReqETH = "https://api.pro.coinbase.com/products/ETH-USD/ticker"

// KrakenReqBTC is the kraken ticker from the API
var KrakenReqBTC = "https://api.kraken.com/0/public/Ticker?pair=BTCUSD"

// KrakenReqETH is the kraken ticker from the API
var KrakenReqETH = "https://api.kraken.com/0/public/Ticker?pair=ETHUSD"

// BinanceVolBTC is the binance ticker from the API
var BinanceVolBTC = "https://api.binance.com/api/v1/ticker/24hr?symbol=BTCUSDT"

// BinanceVolETH is the binance ticker from the API
var BinanceVolETH = "https://api.binance.com/api/v1/ticker/24hr?symbol=ETHUSDT"

// BinanceTickerResponse defines the ticker API response from Binanace
type BinanceTickerResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// BinanceVolumeResponse defines the structure of binance's volume endpoint response
type BinanceVolumeResponse struct {
	// there are other fields as well, but we ignore them for now
	Symbol string `json:"symbol"`
	Volume string `json:"volume"`
}

// CoinbaseTickerResponse defines the structure of coinbase's ticker endpoitt response
type CoinbaseTickerResponse struct {
	TradeId int    `json:"trade_id"`
	Price   string `json:"price"`
	Volume  string `json:"volume"`
}

// KrakenTickerResponse defines the structure of kraken's ticker response
type KrakenTickerResponse struct {
	Error  []string `json:"error"`
	Result struct {
		XXBTZUSD struct {
			// there's some additional info here but we don't require that
			C []string // c = last trade closed array(<price>, <lot volume>),
			V []string // volume array(<today>, <last 24 hours>)
		}
		XETHZUSD struct {
			// there's some additional info here but we don't require that
			C []string // c = last trade closed array(<price>, <lot volume>),
			V []string // volume array(<today>, <last 24 hours>)
		}
	}
}

// BinanceTicker gets price data from Binance
func BinanceTicker(coin string) (float64, error) {
	var data []byte
	var err error
	if coin == "BTC" {
		data, err = erpc.GetRequest(BinanceReqBTC)
	} else if coin == "ETH" {
		data, err = erpc.GetRequest(BinanceReqETH)
	}
	if err != nil {
		log.Println("did not get response", err)
		return -1, errors.Wrap(err, "did not get response from Binance API")
	}

	var response BinanceTickerResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return -1, errors.Wrap(err, "could not unmarshal response")
	}

	if coin == "BTC" && response.Symbol != "BTCUSDT" {
		return -1, errors.New("ticker symbols don't match with API response")
	} else if coin == "ETH" && response.Symbol != "ETHUSDT" {
		return -1, errors.New("ticker symbols don't match with API response")
	}
	// response.Price is in string, need to convert it to float
	price, err := utils.ToFloat(response.Price)
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return math.Round(price*1000) / 1000, nil
}

// BinanceVolume gets volume data from Binance
func BinanceVolume(coin string) (float64, error) {
	var data []byte
	var err error
	if coin == "BTC" {
		data, err = erpc.GetRequest(BinanceVolBTC)
	} else if coin == "ETH" {
		data, err = erpc.GetRequest(BinanceVolETH)
	}
	if err != nil {
		log.Println("did not get response", err)
		return -1, errors.Wrap(err, "did not get response from Binance API")
	}

	var response BinanceVolumeResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return -1, errors.Wrap(err, "could not unmarshal response")
	}

	if coin == "BTC" && response.Symbol != "BTCUSDT" {
		return -1, errors.New("ticker symbols don't match with API response")
	} else if coin == "ETH" && response.Symbol != "ETHUSDT" {
		return -1, errors.New("ticker symbols don't match with API response")
	}

	volume, err := utils.ToFloat(response.Volume)
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return math.Round(volume*1000) / 1000, nil // volume is in BTC and not usd
}

// CoinbaseTicker gets ticker data from coinbase
func CoinbaseTicker(coin string) (float64, error) {
	var data []byte
	var err error
	if coin == "BTC" {
		data, err = erpc.GetRequest(CoinbaseReqBTC)
	} else if coin == "ETH" {
		data, err = erpc.GetRequest(CoinbaseReqETH)
	}
	if err != nil {
		log.Println("did not get response", err)
		return -1, errors.Wrap(err, "did not get response from Coinbase API")
	}

	var response CoinbaseTickerResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return -1, errors.Wrap(err, "could not unmarshal response")
	}

	// response.Price is in string, need to convert it to float
	price, err := utils.ToFloat(response.Price)
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return math.Round(price*1000) / 1000, nil
}

// CoinbaseVolume gets volume data from coinbase
func CoinbaseVolume(coin string) (float64, error) {
	var data []byte
	var err error
	if coin == "BTC" {
		data, err = erpc.GetRequest(CoinbaseReqBTC)
	} else if coin == "ETH" {
		data, err = erpc.GetRequest(CoinbaseReqETH)
	}
	if err != nil {
		log.Println("did not get response", err)
		return -1, errors.Wrap(err, "did not get response from Coinbase API")
	}

	var response CoinbaseTickerResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return -1, errors.Wrap(err, "could not unmarshal response")
	}

	// response.Price is in string, need to convert it to float
	volume, err := utils.ToFloat(response.Volume)
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return math.Round(volume*1000) / 1000, nil
}

// KrakenTicker gets ticker data from kraken
func KrakenTicker(coin string) (float64, error) {
	var data []byte
	var err error
	if coin == "BTC" {
		data, err = erpc.GetRequest(KrakenReqBTC)
	} else if coin == "ETH" {
		data, err = erpc.GetRequest(KrakenReqETH)
	}
	if err != nil {
		log.Println("did not get response", err)
		return -1, errors.Wrap(err, "did not get response from Kraken API")
	}

	var response KrakenTickerResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return -1, errors.Wrap(err, "could not unmarshal response")
	}

	// response.Price is in string, need to convert it to float
	var price float64
	if coin == "BTC" {
		price, err = utils.ToFloat(response.Result.XXBTZUSD.C[0])
	} else if coin == "ETH" {
		price, err = utils.ToFloat(response.Result.XETHZUSD.C[0])
	}

	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return math.Round(price*1000) / 1000, nil
}

// KrakenVolume gets volume data from kraken
func KrakenVolume(coin string) (float64, error) {
	var data []byte
	var err error
	if coin == "BTC" {
		data, err = erpc.GetRequest(KrakenReqBTC)
	} else if coin == "ETH" {
		data, err = erpc.GetRequest(KrakenReqETH)
	}
	if err != nil {
		log.Println("did not get response", err)
		return -1, errors.Wrap(err, "did not get response from Kraken API")
	}

	var response KrakenTickerResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return -1, errors.Wrap(err, "could not unmarshal response")
	}

	var volume float64
	if coin == "BTC" {
		volume, err = utils.ToFloat(response.Result.XXBTZUSD.V[1]) // we want volume over the last 24 hours
	} else if coin == "ETH" {
		volume, err = utils.ToFloat(response.Result.XETHZUSD.V[1]) // we want volume over the last 24 hours
	}
	if err != nil {
		return -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return math.Round(volume*1000) / 1000, nil
}

// Collate collates multiple exchange data
func Collate(coin string) (float64, error) {
	var wg sync.WaitGroup

	binanceVolume, cbVolume, krakenVolume,
		binanceTicker, cbTicker, krakenTicker := 0.0, 0.0, 0.0, 0.0, 0.0, 0.0

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var err error
		binanceVolume, err = CoinbaseVolume(coin)
		if err != nil {
			log.Println(err)
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var err error
		cbVolume, err = CoinbaseVolume(coin)
		if err != nil {
			log.Println(err)
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var err error
		krakenVolume, err = KrakenVolume(coin)
		if err != nil {
			log.Println(err)
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var err error
		binanceTicker, err = CoinbaseTicker(coin)
		if err != nil {
			log.Println(err)
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var err error
		cbTicker, err = CoinbaseTicker(coin)
		if err != nil {
			log.Println(err)
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var err error
		krakenTicker, err = KrakenTicker(coin)
		if err != nil {
			log.Println(err)
		}
	}(&wg)

	wg.Wait()

	netVolume := binanceVolume + cbVolume + krakenVolume

	// return weighted average of all the prices
	return binanceTicker*(binanceVolume/netVolume) +
		cbTicker*(cbVolume/netVolume) +
		krakenTicker*(krakenVolume/netVolume), nil
}
