package internal

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
)

type ResultJson struct {
	SpentMonthly      []GasSpentMonthly                  `json:"spent_monthly"`
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

func (j *ResultJson) WriteJson() {
	file, _ := json.MarshalIndent(j, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0644)
}

func (j *ResultJson) AddGasSpentMonthly(g GasSpentMonthly) {
	j.SpentMonthly = append(j.SpentMonthly, g)
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
