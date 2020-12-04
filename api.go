package main

import (
	"encoding/json"
	"log"
	"math"
	"strings"

	errors "github.com/pkg/errors"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
)

// BinanceReqBTC is the binance ticker from the API
var BinanceReqBTC = "https://api.binance.com/api/v1/ticker/price?symbol=BTCUSDT"
var BinanceReqETH = "https://api.binance.com/api/v1/ticker/price?symbol=ETHUSDT"
var BinanceReqXRP = "https://api.binance.com/api/v1/ticker/price?symbol=XRPUSDT"

var BinanceVolBTC = "https://api.binance.com/api/v1/ticker/24hr?symbol=BTCUSDT"
var BinanceVolETH = "https://api.binance.com/api/v1/ticker/24hr?symbol=ETHUSDT"
var BinanceVolXRP = "https://api.binance.com/api/v1/ticker/24hr?symbol=XRPUSDT"

var CoinbaseReqBTC = "https://api.pro.coinbase.com/products/BTC-USD/ticker"
var CoinbaseReqETH = "https://api.pro.coinbase.com/products/ETH-USD/ticker"
var CoinbaseReqXRP = "https://api.pro.coinbase.com/products/XRP-USD/ticker"

var KrakenReqBTC = "https://api.kraken.com/0/public/Ticker?pair=BTCUSD"
var KrakenReqETH = "https://api.kraken.com/0/public/Ticker?pair=ETHUSD"
var KrakenReqXRP = "https://api.kraken.com/0/public/Ticker?pair=XRPUSD"

var BitfinexReqBTC = "https://api-pub.bitfinex.com/v2/tickers?symbols=tBTCUSD"
var BitfinexReqETH = "https://api-pub.bitfinex.com/v2/tickers?symbols=tETHUSD"
var BitfinexReqXRP = "https://api-pub.bitfinex.com/v2/tickers?symbols=tXRPUSD"

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
		XXRPZUSD struct {
			// there's some additional info here but we don't require that
			C []string // c = last trade closed array(<price>, <lot volume>),
			V []string // volume array(<today>, <last 24 hours>)
		}
	}
}

type BitfinexTickerResponse struct {
	Price  string
	Volume string
}

// BinanceTicker gets price data from Binance
func BinanceTicker(coin string) (float64, error) {
	var data []byte
	var err error
	if coin == "BTC" {
		data, err = erpc.GetRequest(BinanceReqBTC)
	} else if coin == "ETH" {
		data, err = erpc.GetRequest(BinanceReqETH)
	} else if coin == "XRP" {
		data, err = erpc.GetRequest(BinanceReqXRP)
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
	} else if coin == "XRP" && response.Symbol != "XRPUSDT" {
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
	} else if coin == "XRP" {
		data, err = erpc.GetRequest(BinanceVolXRP)
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
func CoinbaseTicker(coin string) (float64, float64, error) {
	var data []byte
	var err error
	if coin == "BTC" {
		data, err = erpc.GetRequest(CoinbaseReqBTC)
	} else if coin == "ETH" {
		data, err = erpc.GetRequest(CoinbaseReqETH)
	} else if coin == "XRP" {
		data, err = erpc.GetRequest(CoinbaseReqXRP)
	}
	if err != nil {
		log.Println("did not get response", err)
		return -1, -1, errors.Wrap(err, "did not get response from Coinbase API")
	}

	var response CoinbaseTickerResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return -1, -1, errors.Wrap(err, "could not unmarshal response")
	}

	// response.Price is in string, need to convert it to float
	price, err := utils.ToFloat(response.Price)
	if err != nil {
		return -1, -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	// response.Price is in string, need to convert it to float
	volume, err := utils.ToFloat(response.Volume)
	if err != nil {
		return -1, -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return math.Round(price*1000) / 1000, math.Round(volume*1000) / 1000, nil
}

// KrakenTicker gets ticker data from kraken
func KrakenTicker(coin string) (float64, float64, error) {
	var data []byte
	var err error
	if coin == "BTC" {
		data, err = erpc.GetRequest(KrakenReqBTC)
	} else if coin == "ETH" {
		data, err = erpc.GetRequest(KrakenReqETH)
	} else if coin == "XRP" {
		data, err = erpc.GetRequest(KrakenReqXRP)
	}
	if err != nil {
		log.Println("did not get response", err)
		return -1, -1, errors.Wrap(err, "did not get response from Kraken API")
	}

	var response KrakenTickerResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return -1, -1, errors.Wrap(err, "could not unmarshal response")
	}

	// response.Price is in string, need to convert it to float
	var price float64
	var volume float64

	var err1 error
	var err2 error

	if coin == "BTC" {
		price, err1 = utils.ToFloat(response.Result.XXBTZUSD.C[0])
		volume, err2 = utils.ToFloat(response.Result.XXBTZUSD.V[1]) // we want volume over the last 24 hours
	} else if coin == "ETH" {
		price, err1 = utils.ToFloat(response.Result.XETHZUSD.C[0])
		volume, err2 = utils.ToFloat(response.Result.XETHZUSD.V[1]) // we want volume over the last 24 hours
	} else if coin == "XRP" {
		price, err1 = utils.ToFloat(response.Result.XXRPZUSD.C[0])
		volume, err2 = utils.ToFloat(response.Result.XXRPZUSD.V[1]) // we want volume over the last 24 hours
	}

	if err1 != nil || err2 != nil {
		return -1, -1, errors.Wrap(err, "could not convert price from string to float, quitting!")
	}

	return math.Round(price*1000) / 1000, math.Round(volume*1000) / 1000, nil
}

// BitfinexTicker gets ticker data from kraken
func BitfinexTicker(coin string) (float64, float64, error) {
	var data []byte
	var err error
	if coin == "BTC" {
		data, err = erpc.GetRequest(BitfinexReqBTC)
	} else if coin == "ETH" {
		data, err = erpc.GetRequest(BitfinexReqETH)
	} else if coin == "XRP" {
		data, err = erpc.GetRequest(BitfinexReqXRP)
	}
	if err != nil {
		log.Println("did not get response", err)
		return -1, -1, errors.Wrap(err, "did not get response from Kraken API")
	}

	response := string(data)
	response = response[1:]
	response = response[0 : len(response)-1]

	responseArr := strings.Split(response, ",")

	price, err := utils.ToFloat(responseArr[1])
	if err != nil {
		return -1, -1, errors.Wrap(err, "did not get response from Kraken API")
	}

	volume, err := utils.ToFloat(responseArr[8])
	if err != nil {
		return -1, -1, errors.Wrap(err, "did not get response from Kraken API")
	}

	// SYMBOL,BID, BID_SIZE, ASK, ASK_SIZE, DAILY_CHANGE, DAILY_CHANGE_RELATIVE, LAST_PRICE, VOLUME, HIGH, LOW
	return math.Round(price*1000) / 1000, math.Round(volume*1000) / 1000, nil
}
