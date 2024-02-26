package converter

import (
	"encoding/json"
	"net/http"
	"strconv"
	"errors"
)

type TickerResponse struct {
	Data []TickerData `json:"data"`
}

type TickerData struct {
	Symbol    string `json:"symbol"`
	Value     string `json:"value"`
	Sources   int    `json:"sources"`
	UpdatedAt string `json:"updated_at"`
}

func BTCEURConverter(amount float64, which string) (float64, error) {
    response, err := http.Get("http://api-cryptopia.adca.sh/v1/prices/ticker?symbol=BTC%2FEUR")
    if err != nil {
        return 0, err
    }
    defer response.Body.Close()

    var tickerResponse TickerResponse
    err = json.NewDecoder(response.Body).Decode(&tickerResponse)
    if err != nil {
        return 0, err
    }

    if len(tickerResponse.Data) == 0 {
        return 0, errors.New("no data found in ticker response")
    }

    valueStr := tickerResponse.Data[0].Value
    value, err := strconv.ParseFloat(valueStr, 64)
    if err != nil {
        return 0, err
    }

    switch which {
    case "BTCTOEUR":
        eurValue := amount * value
        return eurValue, nil
    case "EURTOBTC":
        if value == 0 {
            return 0, errors.New("EUR to BTC conversion rate is zero")
        }
        BTCValue := amount / value
        return BTCValue, nil
    default:
        return 0, errors.New("invalid conversion type")
    }
}

