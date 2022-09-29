package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sort"
	"sync"
	"time"
)

type JsonFileStruct struct {
	Ethereum struct {
		Transactions []struct {
			Time           string  `json:"time"`
			GasPrice       float64 `json:"gasPrice"`
			GasValue       float64 `json:"gasValue"`
			Average        float64 `json:"average"`
			MaxGasPrice    float64 `json:"maxGasPrice"`
			MedianGasPrice float64 `json:"medianGasPrice"`
		} `json:"transactions"`
	} `json:"ethereum"`
}

type GasSumCount struct {
	sum   *big.Float
	count int32
}

func ParseJson(url string, path string) {
	jsonFileStruct := JsonFileStruct{}

	client := &http.Client{Timeout: 10 * time.Second}

	response, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&jsonFileStruct)
	if err != nil {
		log.Fatal(err)
	}

	resultJson := ResultJson{}
	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(4)

	// 1)
	go func() {
		gasSpentPerMonth(&jsonFileStruct, &resultJson)
		wg.Done()
	}()
	// 2)
	go func() {
		averageGasPricePerDay(&jsonFileStruct, &resultJson)
		wg.Done()
	}()
	// 3)
	go func() {
		frequencyDistributionPricePerHour(&jsonFileStruct, &resultJson)
		wg.Done()
	}()

	// 4)
	go func() {
		getEntirePeriodGasPaid(&jsonFileStruct, &resultJson)
		wg.Done()
	}()
	wg.Wait()
	duration := time.Since(start)
	fmt.Println("\nВремя выполнения 4 функций:", duration)

	resultJson.WriteJson(path)
}

//	1) Сколько было потрачено gas помесячно.
func gasSpentPerMonth(jsonFileStruct *JsonFileStruct, resultJson *ResultJson) map[time.Month]*big.Float {
	monthGasSum := make(map[time.Month]*big.Float)
	for _, transaction := range jsonFileStruct.Ethereum.Transactions {
		parse, err := time.Parse("06-01-02 15:04", transaction.Time)

		if err != nil {
			log.Fatal(err)
		}

		if sum, ok := monthGasSum[parse.Month()]; ok {
			sum.Add(sum, new(big.Float).SetFloat64(transaction.GasValue))
		} else {
			monthGasSum[parse.Month()] = new(big.Float).SetFloat64(transaction.GasValue)
		}
	}

	for month, sum := range monthGasSum {
		resultJson.AddGasSpentMonthly(GasSpentMonthly{
			Month:    month.String(),
			GasSpent: sum,
		})
	}
	return monthGasSum
}

//	2) Среднюю цену gas за день.
func averageGasPricePerDay(jsonFileStruct *JsonFileStruct, resultJson *ResultJson) {
	avgGasPricePerDay := make(map[time.Time]GasSumCount)
	for _, transaction := range jsonFileStruct.Ethereum.Transactions {
		parse, err := time.Parse("06-01-02 15:04", transaction.Time)
		if err != nil {
			log.Fatal(err)
		}
		currentDay := time.Date(parse.Year(), parse.Month(), parse.Day(), 0, 0, 0, 0, parse.Location())
		if gasHourCount, ok := avgGasPricePerDay[currentDay]; ok {
			avgGasPricePerDay[currentDay] = GasSumCount{
				sum:   gasHourCount.sum.Add(gasHourCount.sum, new(big.Float).SetFloat64(transaction.MedianGasPrice)),
				count: gasHourCount.count + 1,
			}

		} else {
			avgGasPricePerDay[currentDay] = GasSumCount{
				sum:   new(big.Float).SetFloat64(transaction.MedianGasPrice),
				count: 1,
			}
		}
	}

	finalMap := make(map[time.Time]*big.Float)
	for dateTime, gasSumCount := range avgGasPricePerDay {
		finalMap[dateTime] = gasSumCount.sum.Quo(gasSumCount.sum, new(big.Float).SetInt64(int64(gasSumCount.count)))
	}

	keys := make([]time.Time, 0, len(finalMap))
	for k := range finalMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	for _, k := range keys {
		resultJson.AddAverageDailyPrice(
			AverageDailyPrice{
				Date:        k.Format("06-01-02"),
				AvgGasPrice: finalMap[k],
			},
		)
	}
}

//	3) Частотное распределение цены по часам(за весь период).
func frequencyDistributionPricePerHour(jsonFileStruct *JsonFileStruct, resultJson *ResultJson) {
	frequencyDistribution := make(map[int]GasSumCount)
	for _, transaction := range jsonFileStruct.Ethereum.Transactions {
		parse, err := time.Parse("06-01-02 15:04", transaction.Time)
		if err != nil {
			log.Fatal(err)
		}
		if hourPrice, ok := frequencyDistribution[parse.Hour()]; ok {
			frequencyDistribution[parse.Hour()] = GasSumCount{
				sum:   hourPrice.sum.Add(hourPrice.sum, new(big.Float).SetFloat64(transaction.GasPrice)),
				count: hourPrice.count + 1,
			}
		} else {
			frequencyDistribution[parse.Hour()] = GasSumCount{
				sum:   new(big.Float).SetFloat64(transaction.GasPrice),
				count: 1,
			}
		}
	}

	finalMap := make(map[int]*big.Float)
	for hour, gasSumCount := range frequencyDistribution {
		finalMap[hour] = gasSumCount.sum.Quo(gasSumCount.sum, new(big.Float).SetInt64(int64(gasSumCount.count)))
	}

	keys2 := make([]int, 0, len(finalMap))
	for k := range finalMap {
		keys2 = append(keys2, k)
	}
	sort.Slice(keys2, func(i, j int) bool {
		return keys2[i] < (keys2[j])
	})

	for _, k := range keys2 {
		resultJson.AddPriceFrequencyDistributionByHour(
			PriceFrequencyDistributionByHour{
				Hour:     k,
				GasPrice: finalMap[k],
			},
		)
	}
}

//	4) Сколько заплатили за весь период (gas price * value).
func getEntirePeriodGasPaid(jsonFileStruct *JsonFileStruct, resultJson *ResultJson) {
	sum := new(big.Float)
	for _, transaction := range jsonFileStruct.Ethereum.Transactions {
		gasPrice := new(big.Float).SetFloat64(transaction.GasPrice)
		sum.Add(sum, gasPrice.Mul(gasPrice, new(big.Float).SetFloat64(transaction.GasValue)))
	}

	resultJson.AddEntirePeriodPaid(sum)
}
