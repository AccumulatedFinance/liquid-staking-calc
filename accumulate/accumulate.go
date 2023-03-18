package accumulate

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ybbus/jsonrpc/v3"
)

type AccumulateClient struct {
	API      string
	Client   jsonrpc.RPCClient
	Validate *validator.Validate
}

// NewAccumulateClient constructs the Accumulate client
func NewAccumulateClient(apiURL string, timeout time.Duration) *AccumulateClient {

	c := &AccumulateClient{API: apiURL}

	// init validator
	c.Validate = validator.New()

	// set 5 seconds timeout
	opts := &jsonrpc.RPCClientOpts{}
	opts.HTTPClient = &http.Client{
		Timeout: timeout * time.Second,
	}

	c.Client = jsonrpc.NewClientWithOpts(apiURL, opts)

	return c

}
