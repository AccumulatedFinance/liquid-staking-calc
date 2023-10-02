package accumulated

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type AccumulatedFinanceClient struct {
	API      string
	Client   *http.Client
	Validate *validator.Validate
}

// NewAccumulateClient constructs the Accumulate client
func NewAccumulatedFinanceClient(apiURL string, timeout time.Duration) *AccumulatedFinanceClient {

	c := &AccumulatedFinanceClient{API: apiURL}

	// init validator
	c.Validate = validator.New()

	// set 5 seconds timeout
	c.Client = &http.Client{
		Timeout: timeout * time.Second,
	}

	return c

}
