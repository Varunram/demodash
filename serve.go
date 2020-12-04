package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"text/template"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
)

var (
	// APIError is the error returned if something goes wrong with the API
	APIError = "API error, please try again"

	// RenderError is the error reutnred if something goes wrong while rendering the frontend
	RenderError = "Error while rendering html, please try again"
)

func renderHTML() (string, error) {
	doc, err := ioutil.ReadFile("index.html")
	return string(doc), err
}

type base struct {
	BTCPrice   float64
	ETHPrice   float64
	XRPPrice   float64
	LTCPrice   float64
	LINKPrice  float64
	ADAPrice   float64
	BTCVolume  float64
	ETHVolume  float64
	XRPVolume  float64
	LTCVolume  float64
	LINKVolume float64
	ADAVolume  float64
}

// Return is the structure used to feed data to the frontend
var Return struct {
	Binance  base
	Coinbase base
	Kraken   base
	Bitfinex base
}

func frontend() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		doc, err := renderHTML()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError, APIError)
		}
		templates := template.New("template")
		templates.New("doc").Parse(doc)

		var wg sync.WaitGroup
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Binance.BTCPrice, err = BinanceTicker("BTC")
			if err != nil {
				Return.Binance.BTCPrice = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Binance.ETHPrice, err = BinanceTicker("ETH")
			if err != nil {
				Return.Binance.ETHPrice = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Binance.XRPPrice, err = BinanceTicker("XRP")
			if err != nil {
				Return.Binance.XRPPrice = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Binance.LTCPrice, err = BinanceTicker("LTC")
			if err != nil {
				Return.Binance.LTCPrice = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Binance.LINKPrice, err = BinanceTicker("LINK")
			if err != nil {
				Return.Binance.LINKPrice = -1
			}
		}(&wg)
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Binance.ADAPrice, err = BinanceTicker("ADA")
			if err != nil {
				Return.Binance.ADAPrice = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Binance.BTCVolume, err = BinanceVolume("BTC")
			if err != nil {
				Return.Binance.BTCVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Binance.ETHVolume, err = BinanceVolume("ETH")
			if err != nil {
				Return.Binance.ETHVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Binance.XRPVolume, err = BinanceVolume("XRP")
			if err != nil {
				Return.Binance.XRPVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Binance.LTCVolume, err = BinanceVolume("LTC")
			if err != nil {
				Return.Binance.LTCVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Binance.LINKVolume, err = BinanceVolume("LINK")
			if err != nil {
				Return.Binance.LINKVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Binance.ADAVolume, err = BinanceVolume("ADA")
			if err != nil {
				Return.Binance.ADAVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Coinbase.BTCPrice, Return.Coinbase.BTCVolume, err = CoinbaseTicker("BTC")
			if err != nil {
				Return.Coinbase.BTCPrice = -1
				Return.Coinbase.BTCVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Coinbase.ETHPrice, Return.Coinbase.ETHVolume, err = CoinbaseTicker("ETH")
			if err != nil {
				Return.Coinbase.ETHPrice = -1
				Return.Coinbase.ETHVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Coinbase.XRPPrice, Return.Coinbase.XRPVolume, err = CoinbaseTicker("XRP")
			if err != nil {
				Return.Coinbase.XRPPrice = -1
				Return.Coinbase.XRPVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Coinbase.LTCPrice, Return.Coinbase.LTCVolume, err = CoinbaseTicker("LTC")
			if err != nil {
				Return.Coinbase.LTCPrice = -1
				Return.Coinbase.LTCVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Coinbase.LINKPrice, Return.Coinbase.LINKVolume, err = CoinbaseTicker("LINK")
			if err != nil {
				Return.Coinbase.LINKPrice = -1
				Return.Coinbase.LINKVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Kraken.BTCPrice, Return.Kraken.BTCVolume, err = KrakenTicker("BTC")
			if err != nil {
				Return.Kraken.BTCPrice = -1
				Return.Kraken.BTCVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Kraken.ETHPrice, Return.Kraken.ETHVolume, err = KrakenTicker("ETH")
			if err != nil {
				Return.Kraken.ETHPrice = -1
				Return.Kraken.ETHVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Kraken.XRPPrice, Return.Kraken.XRPVolume, err = KrakenTicker("XRP")
			if err != nil {
				Return.Kraken.XRPPrice = -1
				Return.Kraken.XRPVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Kraken.LTCPrice, Return.Kraken.LTCVolume, err = KrakenTicker("LTC")
			if err != nil {
				Return.Kraken.LTCPrice = -1
				Return.Kraken.LTCVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Kraken.LINKPrice, Return.Kraken.LINKVolume, err = KrakenTicker("LINK")
			if err != nil {
				Return.Kraken.LINKPrice = -1
				Return.Kraken.LINKVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Kraken.ADAPrice, Return.Kraken.ADAVolume, err = KrakenTicker("ADA")
			if err != nil {
				Return.Kraken.ADAPrice = -1
				Return.Kraken.ADAVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Bitfinex.BTCPrice, Return.Bitfinex.BTCVolume, err = BitfinexTicker("BTC")
			if err != nil {
				Return.Bitfinex.BTCPrice = -1
				Return.Bitfinex.BTCVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Bitfinex.ETHPrice, Return.Bitfinex.ETHVolume, err = BitfinexTicker("ETH")
			if err != nil {
				Return.Bitfinex.ETHPrice = -1
				Return.Bitfinex.ETHVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Bitfinex.XRPPrice, Return.Bitfinex.XRPVolume, err = BitfinexTicker("XRP")
			if err != nil {
				Return.Bitfinex.XRPPrice = -1
				Return.Bitfinex.XRPVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Bitfinex.LTCPrice, Return.Bitfinex.LTCVolume, err = BitfinexTicker("LTC")
			if err != nil {
				Return.Bitfinex.LTCPrice = -1
				Return.Bitfinex.LTCVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Bitfinex.LINKPrice, Return.Bitfinex.LINKVolume, err = BitfinexTicker("LINK")
			if err != nil {
				Return.Bitfinex.LINKPrice = -1
				Return.Bitfinex.LINKVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Bitfinex.ADAPrice, Return.Bitfinex.ADAVolume, err = BitfinexTicker("ADA")
			if err != nil {
				Return.Bitfinex.ADAPrice = -1
				Return.Bitfinex.ADAVolume = -1
			}
		}(&wg)

		wg.Wait()
		templates.Lookup("doc").Execute(w, Return)
	})
}

func serveStatic() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
}

func startServer(portx int, insecure bool) {
	frontend()
	serveStatic()

	port, err := utils.ToString(portx)
	if err != nil {
		log.Fatal("Port not string")
	}

	log.Println("Starting RPC Server on Port: ", port)
	if insecure {
		log.Println("starting server in insecure mode")
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(":"+port, "certs/server.crt", "certs/server.key", nil))
	}
}
