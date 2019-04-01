package app

import (
	"github.com/gorilla/mux"
    "github.com/go-pg/pg"
    "github.com/go-pg/pg/orm"
	"handlers"
	"models"
	"net/http"
	"config"
	"utils"
	"fmt"
)

type App struct {
	Router 		*mux.Router
	Config 		*config.Config
	Database	*pg.DB
	RateLimiter	*utils.RateLimiter
}

func (a *App) Initialize(config *config.Config) {
	a.Config = config

	// Specify the list of endpoints bellow
	a.Router = mux.NewRouter()

	a.Router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		handlers.Ping(a.Database, a.Config, w, r)
	}).Methods("GET")

	a.Router.HandleFunc("/price/sync", func(w http.ResponseWriter, r *http.Request) {
		handlers.PriceSync(a.RateLimiter, a.Database, a.Config, w, r)
	}).Methods("GET")

	a.Router.HandleFunc("/price/async", func(w http.ResponseWriter, r *http.Request) {
		handlers.PriceAsync(a.RateLimiter, a.Database, a.Config, w, r)
	}).Methods("GET")

	a.Router.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		handlers.History(a.Database, a.Config, w, r)
	}).Methods("GET")

	a.Router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		handlers.Health(a.Database, a.Config, w, r)
	}).Methods("GET")

	// End endpoints list

	// Connect to the database
	a.Database = pg.Connect(&pg.Options{
        User: 		a.Config.DB.User,
        Password: 	a.Config.DB.Pass,
        Database: 	a.Config.DB.Name,
        Addr:		fmt.Sprintf("%s:%d", a.Config.DB.Host, a.Config.DB.Port),
    })

    e := createSchema(a.Database)
    if e != nil {
    	fmt.Printf(string(e.Error()))
    }

    // Initialize the rate limiter for AlphaVantage
    var r utils.RateLimiter
    a.RateLimiter = &r

    go a.RateLimiter.Initialize(a.Config.DataProvider.RPM)
}

func (a *App) Run() {
	addr := fmt.Sprintf("%s:%d", a.Config.Server.Host, a.Config.Server.Port)

    http.ListenAndServe(addr, a.Router)

    defer a.Database.Close()
}

func createSchema(db *pg.DB) error {
    for _, model := range []interface{}{(*models.PriceModel)(nil)} {
        err := db.CreateTable(model, &orm.CreateTableOptions{
        	IfNotExists: true,
        })

        if err != nil {
            return err
        }
    }
    return nil
}
