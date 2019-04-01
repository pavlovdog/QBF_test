package handlers
import (
    "github.com/go-pg/pg"
	"net/http"
	"config"
)

type HealthSerializer struct {
	Ok	bool	`json:"ok"`
}

func Health(database *pg.DB, config *config.Config, w http.ResponseWriter, r *http.Request) {
	
}

// Works great, as usual
func Ping(database *pg.DB, config *config.Config, w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Pong"))
}