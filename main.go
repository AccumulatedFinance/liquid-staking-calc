package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/AccumulatedFinance/liquid-staking-calc/accumulate"
)

const API_URL = "https://mainnet.accumulatenetwork.io/v2"
const REWARDS_TOKEN_ACCOUNT = "acc://accumulated.acme/staking-rewards"
const TREASURY_TOKEN_ACCOUNT = "acc://accumulated.acme/treasury"
const LIQUID_STAKING_TOKEN_ACCOUNT = "acc://accumulated.acme/staking"
const INCENTIVES_TOKEN_ACCOUNT = "acc://accumulated.acme/incentives"

type Output struct {
	URL          string `json:"url"`
	Amount       int64  `json:"-"`
	AmountString string `json:"amount"`
	Share        int    `json:"-"`
}

type Outputs struct {
	Items []*Output
}

type Balance struct {
	Balance int64 `json:"balance"`
}

// FromString parses balance from string
func (b *Balance) FromString(s string) {
	b.Balance, _ = strconv.ParseInt(s, 10, 64)
}

// String converts balance into human readable format
func (b *Balance) Human() string {
	hr := float64(b.Balance) * math.Pow10(-8)
	return fmt.Sprintf("%.8f", hr)
}

// String converts balance into string
func (b *Balance) String() string {
	return strconv.FormatInt(b.Balance, 10)
}

// FromBalance fills output from balance
func (o *Output) FromBalance(b *Balance) {
	amount := math.Floor(float64(b.Balance) / 10000 * float64(o.Share))
	o.Amount = int64(amount)
	o.AmountString = strconv.FormatInt(o.Amount, 10)
}

// String converts output into human readable format
func (o *Output) String() string {
	hr := float64(o.Amount) * math.Pow10(-8)
	return fmt.Sprintf("%d%% => %s : %.8f ACME", o.Share/100, o.URL, hr)
}

func main() {

	// set distribution
	// https://docs.accumulated.finance/accumulated-finance/fees
	// share is in bps (1% = 100)
	outputs := &Outputs{}
	outputs.Items = append(outputs.Items, &Output{
		URL:   TREASURY_TOKEN_ACCOUNT,
		Share: 2000,
	})
	outputs.Items = append(outputs.Items, &Output{
		URL:   LIQUID_STAKING_TOKEN_ACCOUNT,
		Share: 8000,
	})

	// validate shares
	var totalShare int
	for _, item := range outputs.Items {
		totalShare += item.Share
	}
	if totalShare != 10000 {
		log.Fatal("Expected total shares: ", 10000, ", received: ", totalShare)
	}

	// calculator
	fmt.Println("Calculating liquid staking rewards...")

	client := accumulate.NewAccumulateClient(API_URL, 5)

	fmt.Println("Getting account balance:", REWARDS_TOKEN_ACCOUNT)

	tokenAccount, err := client.QueryTokenAccount(&accumulate.Params{URL: REWARDS_TOKEN_ACCOUNT})
	if err != nil {
		log.Fatal(err)
	}

	balance := &Balance{}
	balance.FromString(tokenAccount.Data.Balance)

	fmt.Println("Balance:", balance.Human(), "ACME")

	for _, item := range outputs.Items {

		item.FromBalance(balance)
		fmt.Println(item)

	}

	fmt.Println("Generating CLI params...")

	json, err := json.Marshal(outputs.Items)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("{ type: sendTokens, to: %v }", string(json))

}