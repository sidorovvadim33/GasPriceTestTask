package internal

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
)

type ResultJson struct {
	GasSpentMonthly []struct {
		Month    string     `json:"month"`
		GasSpent *big.Float `json:"gas_spent"`
	} `json:"spent_monthly"`
	AvgDailyPrice     []AverageDailyPrice                `json:"avg_daily_price"`
	FreqDistribByHour []PriceFrequencyDistributionByHour `json:"freq_distribution_by_hour"`
	EntirePeriodPaid  *big.Float                         `json:"paid"`
}

type GasSpentMonthly struct {
	Month    string     `json:"month"`
	GasSpent *big.Float `json:"gas_spent"`
}

type AverageDailyPrice struct {
	Date        string     `json:"date"`
	AvgGasPrice *big.Float `json:"avg_gas_price"`
}

type PriceFrequencyDistributionByHour struct {
	Hour     int        `json:"hour"`
	GasPrice *big.Float `json:"gas_price"`
}

func (j *ResultJson) WriteJson(path string) {
	file, err := json.MarshalIndent(j, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(path, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (j *ResultJson) AddGasSpentMonthly(g GasSpentMonthly) {
	j.GasSpentMonthly = append(j.GasSpentMonthly, g)
}

func (j *ResultJson) AddAverageDailyPrice(a AverageDailyPrice) {
	j.AvgDailyPrice = append(j.AvgDailyPrice, a)
}

func (j *ResultJson) AddPriceFrequencyDistributionByHour(f PriceFrequencyDistributionByHour) {
	j.FreqDistribByHour = append(j.FreqDistribByHour, f)
}

func (j *ResultJson) AddEntirePeriodPaid(sum *big.Float) {
	j.EntirePeriodPaid = sum
}
