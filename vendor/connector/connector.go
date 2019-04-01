package connector

import (
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"config"
	"strings"
)

type Connector struct {
	Config	*config.DataProvider
}

func (c *Connector) GetPrice(ticker string) (string, error) {
	url := getUrl(c.Config.URL, c.Config.ApiKey, ticker)

	// Make the request
	resp, err := http.Get(url)

	// - Check the request error
	if err != nil {
		return "", err
	}

	// - Convert response
	body, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		return "", err
	}

	escapedBody := escapeDots(string(body))

	price := gjson.Parse(escapedBody).Get("Global Quote").Get("05-price").String()

	return price, nil
}

// Format the API provider endpoint with appropriate parameters, such as API key and so on
func getUrl(basicUrl string, apiKey string, ticker string) string {
	urlWithTicker := strings.Replace(basicUrl, "__SYMBOL__", ticker, 1)
	urlWithKey := strings.Replace(urlWithTicker, "__APIKEY__", apiKey, 1)

	return urlWithKey
}

// Replace the ". " notation with "-"
// E.g. there's a key "05. Price" which becomes the "05-Price"
func escapeDots(r string) string {
	return strings.Replace(r, ". ", "-", -1)
}