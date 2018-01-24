package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"github.com/influxdata/influxdb/client/v2"
)

func writePoints(tags map[string]string, fields map[string]interface{}) {
	fmt.Println(tags)
	clnt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://127.0.0.1:8086",
		Username: "crypto",
		Password: "crypto",
	})
	if err != nil {
		log.Fatal(err)
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "crypto",
		Precision: "us",
	})

	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	pt, err := client.NewPoint(
		"crypto",
		tags,
		fields,
		time.Now(),
	)
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	if err := clnt.Write(bp); err != nil {
		log.Fatal(err)
	}
}

func getValue(curency string, client *bitfinex.Client) float64 {
	change := curency + "usd"
	var price float64
	ticker, err := client.Ticker.Get(change)

	if err != nil {
		fmt.Println("Error")
	} else {
		price, err2 := strconv.ParseFloat(ticker.LastPrice, 64)
		if err2 != nil {
			fmt.Println(err2)
		}
		return price
	}
	return price
}

func getVolume(curency string, client *bitfinex.Client) float64 {

	change := curency + "usd"
	var price float64
	ticker, err := client.Ticker.Get(change)

	if err != nil {
		fmt.Println("Error")
	} else {

		price, err2 := strconv.ParseFloat(ticker.Volume, 64)
		if err2 != nil {
			fmt.Println(err2)
		}
		return price

	}
	return price
}

func main() {
	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")
	
	client := bitfinex.NewClient().Auth(key, secret)

	balance, err := client.Balances.All()

	tags := map[string]string{
		"kind": "value",
	}

	if err != nil {
		fmt.Println(err)
	} else {

		for _, element := range balance {

			amountMap := make(map[string]interface{})
			valueMap := make(map[string]interface{})
			volumeMap := make(map[string]interface{})

			valueMap[element.Currency] = getValue(element.Currency)
			amountMap[element.Currency], err = strconv.ParseFloat(element.Amount, 64)
			volumeMap[element.Currency] = getVolume(element.Currency, client)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(element.Currency, valueMap[element.Currency], amountMap[element.Currency], volumeMap[element.Currency])

			continue // Disable Influx Export
			tags["kind"] = "value"
			writePoints(tags, valueMap)

			tags["kind"] = "amount"
			writePoints(tags, amountMap)

			tags["kind"] = "volume"
			writePoints(tags, volumeMap)

		}
	}
}
