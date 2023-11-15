package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/AccumulatedFinance/liquid-staking-calc/accumulate"
	"github.com/AccumulatedFinance/liquid-staking-calc/accumulated"
)

const ACCUMULATED_FINANCE_API_URL = "https://api.accumulated.finance/v1"

const ACCUMULATE_API_URL = "https://mainnet.accumulatenetwork.io/v2"
const REWARDS_TOKEN_ACCOUNT = "acc://accumulated.acme/staking-rewards"
const TREASURY_TOKEN_ACCOUNT = "acc://accumulated.acme/treasury"
const LIQUID_STAKING_TOKEN_ACCOUNT = "acc://accumulated.acme/staking"
const INCENTIVES_TOKEN_ACCOUNT = "acc://accumulated.acme/incentives"
const PEG_PROTECTION_TOKEN_ACCOUNT = "acc://accumulated.acme/peg-protection"
const WACME_LP_INCENTIVES_TOKEN_ACCOUNT = "acc://accumulated.acme/wacme-lp-incentives"

const STACME_ETHEREUM = "0x7AC168c81F4F3820Fa3F22603ce5864D6aB3C547"
const STACME_ARBITRUM = "0x7AC168c81F4F3820Fa3F22603ce5864D6aB3C547"

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

	// get stACME total supply on Ethereum & Arbitrum
	acfiClient := accumulated.NewAccumulatedFinanceClient(ACCUMULATED_FINANCE_API_URL, 5)

	stacme1, err := acfiClient.GetToken(1, STACME_ETHEREUM)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("stACME Ethereum total supply:", stacme1.TotalSupply)

	stacme42161, err := acfiClient.GetToken(42161, STACME_ARBITRUM)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("stACME Arbitrum total supply:", stacme42161.TotalSupply)

	var share1 float64
	var share42161 float64

	share1 = float64(stacme1.TotalSupply) / (float64(stacme1.TotalSupply) + float64(stacme42161.TotalSupply))
	share42161 = float64(stacme42161.TotalSupply) / (float64(stacme1.TotalSupply) + float64(stacme42161.TotalSupply))

	fmt.Println("stACME Ethereum share:", share1)
	fmt.Println("stACME Arbitrum share:", share42161)

	// set distribution for liquid staking rewards
	// https://docs.accumulated.finance/accumulated-finance/fees
	// share is in bps (1% = 100)
	outputs := &Outputs{}
	outputs.Items = append(outputs.Items, &Output{
		URL:   TREASURY_TOKEN_ACCOUNT,
		Share: 1200,
	})
	outputs.Items = append(outputs.Items, &Output{
		URL:   PEG_PROTECTION_TOKEN_ACCOUNT,
		Share: 800,
	})
	outputs.Items = append(outputs.Items, &Output{
		URL:   LIQUID_STAKING_TOKEN_ACCOUNT,
		Share: int(8000 * share1),
	})
	outputs.Items = append(outputs.Items, &Output{
		URL:   LIQUID_STAKING_TOKEN_ACCOUNT,
		Share: 8000 - int(8000*share1),
	})

	// validate shares
	var totalShare int
	for _, item := range outputs.Items {
		totalShare += item.Share
	}
	if totalShare != 10000 {
		log.Fatal("Expected total shares: ", 10000, ", received: ", totalShare)
	}

	// set distribution for wacme lp incentives
	// https://docs.accumulated.finance/accumulated-finance/fees
	// share is in bps (1% = 100)
	/*
		outputs_wacme := &Outputs{}
		outputs_wacme.Items = append(outputs_wacme.Items, &Output{
			URL:   TREASURY_TOKEN_ACCOUNT,
			Share: 800,
		})
		outputs_wacme.Items = append(outputs_wacme.Items, &Output{
			URL:   INCENTIVES_TOKEN_ACCOUNT,
			Share: 9200,
		})
	*/

	// validate shares
	/*
		var totalShare_wacme int
		for _, item := range outputs_wacme.Items {
			totalShare_wacme += item.Share
		}
		if totalShare_wacme != 10000 {
			log.Fatal("Expected total shares: ", 10000, ", received: ", totalShare_wacme)
		}
	*/

	client := accumulate.NewAccumulateClient(ACCUMULATE_API_URL, 5)

	// liquid staking calculator
	fmt.Println("Calculating liquid staking rewards...")
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

	jsonPayload, err := json.Marshal(outputs.Items)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("'{ type: sendTokens, to: %v }'", string(jsonPayload))
	fmt.Println("")

	// wacme lp incentives calculator
	// removed due to change in WACME/frxETH Curve incentives mechanism
	/*
		fmt.Println("Calculating WACME LP incentives...")
		fmt.Println("Getting account balance:", WACME_LP_INCENTIVES_TOKEN_ACCOUNT)

		tokenAccount, err = client.QueryTokenAccount(&accumulate.Params{URL: WACME_LP_INCENTIVES_TOKEN_ACCOUNT})
		if err != nil {
			log.Fatal(err)
		}

		balance = &Balance{}
		balance.FromString(tokenAccount.Data.Balance)

		fmt.Println("Balance:", balance.Human(), "ACME")

		for _, item := range outputs_wacme.Items {

			item.FromBalance(balance)
			fmt.Println(item)

		}

		fmt.Println("Generating CLI params...")

		jsonPayload, err = json.Marshal(outputs_wacme.Items)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("'{ type: sendTokens, to: %v }'", string(jsonPayload))
		fmt.Println("")
	*/

}
