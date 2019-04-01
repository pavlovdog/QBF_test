package config

type Config struct {
	DataProvider 	*DataProvider
	Logging			*Logging
	Server			*Server
	DB				*DB
}

type DB struct {
	User	string
	Pass	string
	Host	string
	Name	string
	Port	uint
}

type DataProvider struct {
	RPM			int
	ApiKey 		string
	URL 		string
}

type Logging struct {
	Level		string
	SentryKey	string
}

type Server struct {
	Port		uint
	Host		string
}

// TODO: add importing the options from the env for better Docker-compatability
func GetConfig() *Config {
	return &Config {
		DataProvider: &DataProvider {
			RPM: 5,
			ApiKey: "HU8YVUJYWMVOLRII",
			URL: "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=__SYMBOL__&apikey=__APIKEY__",
		},
		Logging: &Logging {
			Level: "info",
			SentryKey: "",
		},
		Server: &Server {
			Port: 8080,
			Host: "127.0.0.1",
		},
		DB: &DB {
			User: "postgres",
			Pass: "123123",
			Port: 5432,
			Host: "localhost",
			Name: "qbf",
		},
	}
}
