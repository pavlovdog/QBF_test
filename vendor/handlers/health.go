package handlers
import (
	"github.com/jinzhu/gorm"
	"net/http"
	"config"
)

type HealthSerializer struct {
	Ok	bool	`json:"ok"`
}

func Health(database *gorm.DB, config *config.Config, w http.ResponseWriter, r *http.Request) {
	
}

// Works great, as usual
func Ping(database *gorm.DB, config *config.Config, w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Pong"))
}