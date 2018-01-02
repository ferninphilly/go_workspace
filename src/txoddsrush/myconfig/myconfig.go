//Package myconfig loads config data into a cfg struct to be used throughout the application
package myconfig

import (
	"encoding/json"
	"log"
	"os"
)

//CONFIG is the Name of our Config file
const CONFIG = "./myconfig/config.json"

//Config struct is the struct to carry our configurations for the application. We will unmarshal from the
//config.json file
type Config struct {
	APIData struct {
		Username   string            `json:"Username"`
		Password   string            `json:"Password"`
		BaseURL    string            `json:"BaseUrl"`
		EndURL     string            `json:"EndUrl"`
		FeedType   string            `json:"FeedType"`
		URLOptions map[string]string `json "UrlOptions`
	} `json:"ApiData"`
	DBData struct {
		Host       string `json:"Host"`
		Port       string `json:"Port"`
		DBUser     string `json:"DBUser"`
		DBPassword string `json:"DBPassword"`
		DBName     string `json:"DBName"`
	} `json:"DBData"`
}

//Generic Error Handler
func HandleError(e error) {
	if e != nil {
		log.Print(e)
	}
}

//LoadConfig will read the config.json file and load it
func (cfg *Config) LoadConfig(file string) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	HandleError(err)
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&cfg)
	HandleError(err)
}

//ReturnConfig loads the config file and returns it
func ReturnConfig() Config {
	var cfg Config
	cfg.LoadConfig(CONFIG)
	return cfg
}
