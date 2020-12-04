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
	BTCPrice  float64
	ETHPrice  float64
	BTCVolume float64
	ETHVolume float64
}

// Return is the structure used to feed data to the frontend
var Return struct {
	Binance  base
	Coinbase base
	Kraken   base
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
			Return.Binance.BTCVolume, err = BinanceVolume("BTC")
			if err != nil {
				Return.Binance.BTCVolume = -1
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
			Return.Binance.ETHVolume, err = BinanceVolume("ETH")
			if err != nil {
				Return.Binance.ETHVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Coinbase.BTCPrice, err = CoinbaseTicker("BTC")
			if err != nil {
				Return.Coinbase.BTCPrice = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Coinbase.BTCVolume, err = CoinbaseVolume("BTC")
			if err != nil {
				Return.Coinbase.BTCVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Coinbase.ETHPrice, err = CoinbaseTicker("ETH")
			if err != nil {
				Return.Coinbase.ETHPrice = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Coinbase.ETHVolume, err = CoinbaseVolume("ETH")
			if err != nil {
				Return.Coinbase.ETHVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Kraken.BTCPrice, err = KrakenTicker("BTC")
			if err != nil {
				Return.Kraken.BTCPrice = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Kraken.BTCVolume, err = KrakenVolume("BTC")
			if err != nil {
				Return.Kraken.BTCVolume = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Kraken.ETHPrice, err = KrakenTicker("ETH")
			if err != nil {
				Return.Kraken.ETHPrice = -1
			}
		}(&wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.Kraken.ETHVolume, err = KrakenVolume("ETH")
			if err != nil {
				Return.Kraken.ETHVolume = -1
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
