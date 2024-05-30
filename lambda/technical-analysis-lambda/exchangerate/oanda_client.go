package exchangerate

import (
	"os"
	"time"

	oanda "github.com/rollend/goanda"
)

type OandaClient interface {
	OandaForexCandles(symbol string, resolution oanda.Granularity, from time.Time, to time.Time) (oanda.InstrumentHistory, error)
}

type OandaForexApi struct {
	client *oanda.Connection
	err    error
}

func (fa *OandaForexApi) OandaForexCandles(symbol string, resolution oanda.Granularity, from time.Time, to time.Time) (oanda.InstrumentHistory, error) {
	return fa.client.GetTimeToCandles(symbol, 500, resolution, from)
}

func NewOandaForexApi() *OandaForexApi {
	accountID := os.Getenv("OANDA_ACCOUNT_ID")
	key := os.Getenv("OANDA_API_KEY")
	config := oanda.ConnectionConfig{
		UserAgent: "v20-golang/0.0.1",
		Timeout:   time.Second * 5,
		Live:      true,
	}
	client, err := oanda.NewConnection(accountID, key, &config)
	return &OandaForexApi{client: client, err: err}
}
