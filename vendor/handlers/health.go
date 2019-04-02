package handlers
import (
	// "github.com/op/go-logging"
	"github.com/jinzhu/gorm"
	"encoding/json"
	"net/http"
	"config"
	"models"
)

// For health check I'll simply use the amount of records
// In the database - no big deal
type HealthSerializer struct {
	Records 	int		`json:"records"`
	Ok			bool	`json:"ok"`
	Msg			string	`json:"msg"`
}

// Request the amount of records in the "price_models" table and return it in a form
// of HealthSerializer
func Health(database *gorm.DB, config *config.Config, w http.ResponseWriter, r *http.Request) {
	var records int
	database.Model(&models.PriceModel{}).Where("").Count(&records)

	response := HealthSerializer{records, true, ""}
    json.NewEncoder(w).Encode(response)	
}

// Works great, as usual
func Ping(database *gorm.DB, config *config.Config, w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Pong"))
}