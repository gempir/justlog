package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/op/go-logging"
)

var (
	cfg config
	// Log logger from go-logging
	Log logging.Logger
)

type config struct {
	IrcAddress string `json:"irc_address"`
	BrokerPass string `json:"broker_pass"`
	APIPort    string `json:"api_port"`
}

func main() {
	Log = initLogger()
	var err error
	cfg, err = readConfig("config.json")
	if err != nil {
		Log.Fatal(err)
	}

}

func initLogger() logging.Logger {
	var logger *logging.Logger
	logger = logging.MustGetLogger("gempbotgo")
	backend := logging.NewLogBackend(os.Stdout, "", 0)

	format := logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000} %{shortfile:-15s} %{level:.4s}%{color:reset} %{message}`,
	)
	logging.SetFormatter(format)
	backendLeveled := logging.AddModuleLevel(backend)
	logging.SetBackend(backendLeveled)
	return *logger
}

func readConfig(path string) (config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	return unmarshalConfig(file)
}

func unmarshalConfig(file []byte) (config, error) {
	err := json.Unmarshal(file, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}