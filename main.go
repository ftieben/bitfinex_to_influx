package main

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"os"
	"strconv"
	"github.com/influxdata/influxdb/client/v2"
	"time"
	"math/rand"
	"log"
)


func writePoints(tags map[string]string, fields map[string]interface{}) {
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



func getValue( curency string) float64 {

	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")

	client := bitfinex.NewClient().Auth(key, secret)
	change := curency + "usd"
	var price float64
	ticker, err := 	client.Ticker.Get(change)



	if err != nil {
		fmt.Println("Error")
	}else {


		price, err2 := strconv.ParseFloat(ticker.LastPrice,64)
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

	tags1 := map[string]string{
		"kind":    "value",
	}
	tags2 := map[string]string{
		"kind":    "amount",
	}


	if err != nil {
		fmt.Println(err)
	}else {


		for _, element := range balance {

			amountMap := make(map[string]interface{})
			valueMap := make(map[string]interface{})

			valueMap[element.Currency] = getValue(element.Currency)
			amountMap[element.Currency], err = strconv.ParseFloat(element.Amount,64)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(element.Currency, valueMap[element.Currency], amountMap[element.Currency])
			//tags["kind"] = "value"
			go writePoints(tags1, valueMap)
			//tags["kind"] = "amount"
			go writePoints(tags2, amountMap)

		}
	}
}
