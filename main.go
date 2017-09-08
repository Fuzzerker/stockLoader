package main

import (
	"fmt"
	"net/http"

	"github.com/gocarina/gocsv"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const apiKey = "62DL09RVRT14MXXN"
const nasdaq = "http://www.nasdaq.com/screening/companies-by-name.aspx?letter=0&exchange=nasdaq&render=download"
const nyse = "http://www.nasdaq.com/screening/companies-by-name.aspx?letter=0&exchange=nyse&render=download"
const amex = "http://www.nasdaq.com/screening/companies-by-name.aspx?letter=0&exchange=amex&render=download"

func main() {

	fmt.Println("Helloadfworld!")

	nasdaqStocks := getStocksFromUrl(nasdaq, "Nasdaq")
	nyseStocks := getStocksFromUrl(nyse, "Nyse")
	amexStocks := getStocksFromUrl(amex, "Amex")
	putInMongo(nasdaqStocks)
	putInMongo(nyseStocks)
	putInMongo(amexStocks)

	fmt.Println(nasdaqStocks[0].Source)
}

func getStocksFromUrl(url string, source string) []Stock {
	res, err := http.Get(nasdaq)
	if err != nil {
		panic(err.Error())
	}
	var stocks []Stock
	err = gocsv.Unmarshal(res.Body, &stocks)

	if err != nil {
		panic(err.Error())
	}

	for i, stock := range stocks {
		stock.Source = source
		stocks[i] = stock
	}

	return stocks
}

func putInMongo(stocks []Stock) {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("Stocks").C("Stocks")

	c.Insert(stocks)

	for _, stock := range stocks {
		_, err = c.Upsert(bson.M{"symbol": stock.Symbol}, stock)
		if err != nil {
			panic("error upserting: " + err.Error())
		}
	}

}

type Stock struct { // Our example struct, you can use "-" to ignore a field
	Symbol       string
	Name         string
	LastSale     string
	MarketCap    string
	IPOyear      string
	Sector       string
	Industry     string `csv:"industry"`
	SummaryQuote string `csv:"Summary Quote"`
	Source       string `csv:"-"`
}
