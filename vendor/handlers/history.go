package handlers
import (
	"github.com/jinzhu/gorm"
	"encoding/json"
	"net/http"
	"config"
	"models"
	// "fmt"
)

type HistorySerializer struct {
	Ok		bool					`json:"ok"`
	Msg		string 					`json:"msg"`
	History []models.PriceModel		`json:"history"`
}

func History(database *gorm.DB, config *config.Config, w http.ResponseWriter, r *http.Request) {
	var h []models.PriceModel
	database.Find(&h)

	response := HistorySerializer{true, "", h}
    json.NewEncoder(w).Encode(response)
}