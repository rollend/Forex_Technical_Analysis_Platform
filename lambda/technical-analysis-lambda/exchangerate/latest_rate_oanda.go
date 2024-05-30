package exchangerate

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"technical-analysis-lambda/finance"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
	oanda "github.com/rollend/goanda"
)

func CreateNewSymbolRate_Oanda(fc FinnhubClient, oc OandaClient, symbols *[4]string, startSeconds int64, endSeconds int64, period string, headline *string) (*[]finance.FinancialDataItem, *[]finance.NewsItem, error) {

	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: os.Getenv("PROJECT_API_KEY"),
	})
	symbolRateItems := []finance.FinancialDataItem{}
	currentTime := time.Now().Unix()

	g := oanda.Granularity(time.Minute)

	for _, symbol := range *symbols {
		// Creates symbol name
		symbolFinnhub := symbol[:3] + "_" + symbol[3:]

		// Gets Forex candles between start and end times
		// Bug exists in Finnhub API where some data outside of time interval is retrieved
		forexCandles, err := oc.OandaForexCandles(symbolFinnhub, g, time.Now(), time.Now())
		if err != nil {
			log.Println("Issue connecting to forex client")
			log.Println(err)
			return nil, nil, errors.New("forex client connection error")
		}
		if forexCandles.Instrument == "" {
			log.Println("No data received")
			return nil, nil, nil
		}

		rateLength := len(forexCandles.Candles)
		for i := 0; i < rateLength; i++ {
			intTimestamp := int64(forexCandles.Candles[i].Time.Unix())
			// Ensure that the retrieved data fits within the specified time interval
			if (intTimestamp > currentTime-startSeconds) && (intTimestamp < currentTime-endSeconds) {
				if i != rateLength-1 && intTimestamp == int64(forexCandles.Candles[i+1].Time.Unix()) {
					continue
				}
				open := float32(forexCandles.Candles[i].Mid.Open)
				high := float32(forexCandles.Candles[i].Mid.High)
				low := float32(forexCandles.Candles[i].Mid.Low)
				close := float32(forexCandles.Candles[i].Mid.Close)
				volume := float32(forexCandles.Candles[i].Volume)
				timestamp := time.Unix(intTimestamp, 0).UTC()

				newSymbolData := finance.SymbolData{Open: open, High: high, Low: low, Close: close, Volume: volume}
				formattedDate := timestamp.Format("2006-01-02")
				formattedTimestamp := timestamp.Format("15:04:05")

				// Check whether there are duplicates in symbolRateItems
				if searchTimeIndex(&symbolRateItems, formattedDate, formattedTimestamp) == -1 {
					switch symbol {
					case "EURUSD":
						newSymbolRateItem := finance.FinancialDataItem{Date: formattedDate, Timestamp: formattedTimestamp, EURUSD: newSymbolData}
						symbolRateItems = append(symbolRateItems, newSymbolRateItem)
					case "GBPUSD":
						newSymbolRateItem := finance.FinancialDataItem{Date: formattedDate, Timestamp: formattedTimestamp, GBPUSD: newSymbolData}
						symbolRateItems = append(symbolRateItems, newSymbolRateItem)
					case "USDJPY":
						newSymbolRateItem := finance.FinancialDataItem{Date: formattedDate, Timestamp: formattedTimestamp, USDJPY: newSymbolData}
						symbolRateItems = append(symbolRateItems, newSymbolRateItem)
					case "AUDCAD":
						newSymbolRateItem := finance.FinancialDataItem{Date: formattedDate, Timestamp: formattedTimestamp, AUDCAD: newSymbolData}
						symbolRateItems = append(symbolRateItems, newSymbolRateItem)
					}
				} else {
					// If Finnhub provides duplicate data, update with the newer duplicate
					index := searchTimeIndex(&symbolRateItems, formattedDate, formattedTimestamp)
					switch symbol {
					case "EURUSD":
						symbolRateItems[index].EURUSD = newSymbolData
					case "GBPUSD":
						symbolRateItems[index].GBPUSD = newSymbolData
					case "USDJPY":
						symbolRateItems[index].USDJPY = newSymbolData
					case "AUDCAD":
						symbolRateItems[index].AUDCAD = newSymbolData
					}
				}
			}
		}
	}

	// Searches for latest market news from Finnhub
	latestNews, err := getMarketNews(fc, &auth, headline)
	if err != nil {
		log.Println(err)
		return &symbolRateItems, nil, nil
	}
	return &symbolRateItems, latestNews, nil
}
