package internal

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

var Config Configuration

const (
	ApiKey           = "6E588445A8"
	PhosconUrl       = "https://phoscon.de/discover"
	CsvPath          = "sensors.csv"
	DelayInSecond    = 30
	TraceHttp        = true
	HttpPort         = ":9000"
	ConfigFileName   = "go-conz-config"
	PathToConfigFile = "."
)

type Configuration struct {
	ApiKey        string
	PhosconUrl    string
	CsvPath       string
	DelayInSecond time.Duration
	TraceHttp     bool
	HttpPort      string
}

func InitConfig() {
	log.Println("Init config...")
	viper.SetDefault("apiKey", ApiKey)
	viper.SetDefault("phosconUrl", PhosconUrl)
	viper.SetDefault("csvPath", CsvPath)
	viper.SetDefault("delayInSecond", DelayInSecond)
	viper.SetDefault("traceHttp", TraceHttp)
	viper.SetDefault("httpPort", HttpPort)
	viper.AutomaticEnv()

	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(PathToConfigFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, get default config: %s\n", err)
	}
	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	log.Println("...Init config finished !")
}
