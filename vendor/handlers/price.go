package handlers
import (
	"github.com/jinzhu/gorm"
	"encoding/json"
	"net/http"
	"connector"
	"strconv"
	"models"
	"config"
	"utils"
	// "fmt"
)

type PriceSerializer struct {
	Price 		string	`json:"price"`
	Ok			bool	`json:"ok"`
	Msg			string	`json:"msg"`
}

// Return the price for ticker in a sync mode
func PriceSync(rl *utils.RateLimiter, database *gorm.DB, config *config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ticker := r.URL.Query().Get("ticker")

	// Initialize the connector object (in this case it implements the AlphaVantage API)
	c := connector.Connector{config.DataProvider}

	rl.WaitForRateLimitApprove()

	// - Request the price
	price, err := c.GetPrice(ticker)

	// Something went wrong, return the error to the user
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := PriceSerializer{"", false, string(err.Error())}
	    json.NewEncoder(w).Encode(response)
		return
	}

	response := PriceSerializer{price, true, ""}
    json.NewEncoder(w).Encode(response)

    // Save price into the db
    go func() {
    	floatPrice, _ := strconv.ParseFloat(price, 64)
    	priceModel := models.PriceModel{Price: floatPrice}
    	priceModel.SavePrice(database)
    }()
}

// Return the price for ticker in an async mode
func PriceAsync(rl *utils.RateLimiter, database *gorm.DB, config *config.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Return some dummy response
	response := PriceSerializer{"", true, "Processing the request"}
    json.NewEncoder(w).Encode(response)

	// Load the price & save it in async mode
	go func() {
		ticker := r.URL.Query().Get("ticker")
		c := connector.Connector{config.DataProvider}

		rl.WaitForRateLimitApprove()

		price, _ := c.GetPrice(ticker)

    	floatPrice, _ := strconv.ParseFloat(price, 64)
    	priceModel := models.PriceModel{Price: floatPrice}
    	priceModel.SavePrice(database)
	}()
}
