package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"github.com/influxdata/influxdb/client/v2"
)

func writePoints(fields map[string]interface{}) {

	// Conntect to Influx
	InfluxHost := os.Getenv("INFLUX_HOST")
	InfluxDB := os.Getenv("INFLUX_DB")
	InfluxUser := os.Getenv("INFLUX_USER")
	InfluxPass := os.Getenv("INFLUX_Pass")

	clnt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     InfluxHost,
		Username: InfluxUser,
		Password: InfluxPass,
	})
	if err != nil {
		fmt.Println(err) //log.Fatal(err)

	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  InfluxDB,
		Precision: "us",
	})

	if err != nil {
		fmt.Println(err) //log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	tags := make(map[string]string)

	pt, err := client.NewPoint(
		"crypto",
		tags,
		fields,
		time.Now(),
	)
	if err != nil {
		fmt.Println(err) //log.Fatal(err)
	}
	bp.AddPoint(pt)

	if err := clnt.Write(bp); err != nil {
		fmt.Println(err) //log.Fatal(err)
	}
}

func getData(element bitfinex.WalletBalance, client *bitfinex.Client) map[string]interface{} {
	change := element.Currency + "usd"
	Ticker, err := client.Ticker.Get(change)

	if err != nil {
		log.Fatal(err)
		return make(map[string]interface{})

	} else {
		returnMap := make(map[string]interface{})
		returnMap["Curency"] = strings.ToUpper(element.Currency)
		returnMap["Volume"] = Ticker.Volume
		returnMap["Ask"] = Ticker.Ask
		returnMap["Bid"] = Ticker.Bid
		returnMap["High"] = Ticker.High
		returnMap["LastPrice"] = Ticker.LastPrice
		returnMap["Low"] = Ticker.Low
		returnMap["Mid"] = Ticker.Mid
		returnMap["Balance"] = element.Amount
		return returnMap
	}
}

func main() {
	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")

	client := bitfinex.NewClient().Auth(key, secret)

	balance, err := client.Balances.All()

	if err != nil {
		log.Fatal(err)
	} else {

		for _, element := range balance {
			fmt.Println(getData(element, client))

			writePoints(getData(element, client))
		}
	}
}
