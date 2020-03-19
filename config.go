package proliferate

import (
	"encoding/json"
	"os"
)

// Config root config populated by config.json
type Config struct {
	//Node     ConfigNode    `json:"node"`
	Logging  Logging       `json:"logging"`
	Couch    CouchConfig   `json:"couchDB"`
	Network  NetworkConfig `json:"network"`
	Instance Instance      `json:"node"`
	Build    BuildConfig
}

// CouchConfig couchConfig populated by config.json
type CouchConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Database string `json:"database"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	State    CouchState
}

// Logging logging read from config.json
type Logging struct {
	Enabled      bool   `json:"enabled"`
	Level        int    `json:"level"`
	Console      bool   `json:"console"`
	File         bool   `json:"file"`
	FileLocation string `json:"fileLocation"`
}

// NetworkConfig networkConfig read from config.json
type NetworkConfig struct {
	Algorithm string   `json:"consensusAlgorithm"`
	Role      int      `json:"role"`
	MaxPeers  int      `json:"maxPeers"`
	Discovery []string `json:"DiscoveryURL"`
}

// Instance node specific settings from config.json
type Instance struct {
	MemoryRecordLimit int `json:"memoryRecordLimit"`
}

// BuildConfig build specific settings from config.json
type BuildConfig struct {
	IdentityFile   string
	IdentityFolder string
	ConfigFile     string
	CertFile       string
	CertExpYears   int
	KeyFile        string
}

var buildConfig = BuildConfig{
	IdentityFolder: ".id",
	IdentityFile:   "id.json",
	ConfigFile:     "config.json",
	CertFile:       "id.cert",
	CertExpYears:   5,
	KeyFile:        "id.pem",
}

// LoadConfig returns json as struct
func LoadConfig() Config {
	var config Config

	file := buildConfig.ConfigFile

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

	config.Build = buildConfig

	return (config)
}
