package accumulated

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type Token struct {
	Address     string     `json:"address" validate:"required,eth_addr"`
	Symbol      string     `json:"symbol" validate:"required"`
	Decimals    int64      `json:"decimals"`
	TotalSupply int64      `json:"totalSupply"`
	ChainID     int        `json:"chainId"`
	Price       float64    `json:"price"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

type TokenResponse struct {
	Result *Token `json:"result"`
}

// QueryADI gets ADI info
func (c *AccumulatedFinanceClient) GetToken(chainID int, address string) (*Token, error) {

	token := &TokenResponse{}

	url := c.API + filepath.Join("/tokens", strconv.Itoa(chainID), address)

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}

	if err := c.Validate.Struct(token.Result); err != nil {
		return nil, err
	}

	return token.Result, nil

}
