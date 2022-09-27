package main

import (
	"fmt"
	"gasPriceTestTask/internal"
	"time"
)

func main() {
	start := time.Now()
	url := "https://raw.githubusercontent.com/CryptoRStar/GasPriceTestTask/main/gas_price.json"

	internal.ParseJson(url)

	duration := time.Since(start)

	fmt.Println("\nВремя выполнения:", duration)
}
