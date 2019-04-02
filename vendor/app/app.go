package app

import (
	"github.com/op/go-logging"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"handlers"
	"models"
	"net/http"
	"config"
	"utils"
	"fmt"
	"os"
)

type App struct {
	Router 		*mux.Router
	Config 		*config.Config
	Database	*gorm.DB
	RateLimiter	*utils.RateLimiter
	Logger		*logging.Logger
}

func (a *App) Initialize(config *config.Config) {
	a.Config = config

	a.setupLogging()
	a.setupRouter()
	a.setupDatabase()
	a.setupRateLimiter()
}

func (a *App) setupLogging() {
    a.Logger = logging.MustGetLogger("example")
    var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{level} %{id:03x}%{color:reset} %{message}`,
	)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	level := logging.AddModuleLevel(backend)
	level.SetLevel(logging.ERROR, "")
	logging.SetBackend(level, backendFormatter)
}

func(a *App) logHTTPRequest(URL string, IP string) {
	a.Logger.Debug(fmt.Sprintf("Request for %s from %s", URL, IP))
}

func (a *App) setupRouter() {
	// Specify the list of endpoints bellow
	a.Router = mux.NewRouter()

	a.Router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		a.logHTTPRequest("/ping", r.RemoteAddr)
		handlers.Ping(a.Database, a.Config, w, r)
	}).Methods("GET")

	a.Router.HandleFunc("/price/sync", func(w http.ResponseWriter, r *http.Request) {
		a.logHTTPRequest("/price/sync", r.RemoteAddr)
		handlers.PriceSync(a.RateLimiter, a.Database, a.Config, w, r)
	}).Methods("GET")

	a.Router.HandleFunc("/price/async", func(w http.ResponseWriter, r *http.Request) {
		a.logHTTPRequest("/price/async", r.RemoteAddr)
		handlers.PriceAsync(a.RateLimiter, a.Database, a.Config, w, r)
	}).Methods("GET")

	a.Router.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		a.logHTTPRequest("/history", r.RemoteAddr)
		handlers.History(a.Database, a.Config, w, r)
	}).Methods("GET")

	a.Router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		a.logHTTPRequest("/health", r.RemoteAddr)
		handlers.Health(a.Database, a.Config, w, r)
	}).Methods("GET")
}

func (a *App) setupDatabase() {
	a.Logger.Debug("Using database", a.Config.DB.Name)

	db, err := gorm.Open("sqlite3", a.Config.DB.Name)
	if err != nil {
		a.Logger.Error("Failed to use database", a.Config.DB.Name)
		panic("")
	}

	a.Database = db
	a.Database.AutoMigrate(&models.PriceModel{})
}

func (a *App) setupRateLimiter() {
	a.Logger.Debug("Using rate limiter with RPM", a.Config.DataProvider.RPM)

    var r utils.RateLimiter
    a.RateLimiter = &r

    go a.RateLimiter.Initialize(a.Config.DataProvider.RPM)
}

func (a *App) Run() {
	addr := fmt.Sprintf("%s:%d", a.Config.Server.Host, a.Config.Server.Port)

	a.Logger.Info("Starting server on", addr)

    http.ListenAndServe(addr, a.Router)

    defer a.Database.Close()
}
