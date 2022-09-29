package main

import (
	"flag"
	"fmt"
	"gasPriceTestTask/internal"
	"time"
)

func main() {
	start := time.Now()
	url := "https://raw.githubusercontent.com/CryptoRStar/GasPriceTestTask/main/gas_price.json"

	var jsonFilePath string
	flag.StringVar(&jsonFilePath, "json-path", "test.json", "defines json file path")
	flag.Parse()

	internal.ParseJson(url, jsonFilePath)

	duration := time.Since(start)
	fmt.Println("\nВремя выполнения всей программы:", duration)
}
