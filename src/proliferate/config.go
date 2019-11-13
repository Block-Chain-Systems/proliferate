package proliferate

import (
	"encoding/json"
	"os"
)

type Config struct {
	Logging Logging     `json:"logging"`
	Couch   CouchConfig `json:"couchDB"`
}

type CouchConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

type Logging struct {
	Enabled      bool   `json:"enabled"`
	Level        int    `json:"level"`
	Console      bool   `json:"console"`
	File         bool   `json:"file"`
	FileLocation string `json:"fileLocation"`
}

// LoadJSON returns json as struct (TODO!)
func LoadConfig() Config {
	var config Config

	file := "config.json"

	configFile, err := os.Open(file)
	defer configFile.Close()

	if err != nil {
		LogRaw(Message{
			Level: 1,
			Text:  err.Error(),
		})
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return (config)
}
