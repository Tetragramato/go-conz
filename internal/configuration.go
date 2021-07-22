package internal

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

var Config *Configuration

const (
	PhosconUrl       = "https://phoscon.de/discover"
	DatabasePath     = "./goconz-sensors"
	DelayInSecond    = 30
	TraceHttp        = true
	HttpPort         = ":9000"
	ConfigFileName   = "./go-conz-config"
	PathToConfigFile = "."
	EnvPrefix        = "GOCONZ"
	ReadOnly         = false
)

type Configuration struct {
	ApiKey        string
	PhosconUrl    string
	DatabasePath  string
	DelayInSecond time.Duration
	TraceHttp     bool
	HttpPort      string
	ReadOnly      bool
}

func InitConfig() {
	log.Println("Init config...")
	viper.SetDefault("phosconUrl", PhosconUrl)
	viper.SetDefault("databasePath", DatabasePath)
	viper.SetDefault("delayInSecond", DelayInSecond)
	viper.SetDefault("traceHttp", TraceHttp)
	viper.SetDefault("httpPort", HttpPort)
	viper.SetDefault("readOnly", ReadOnly)
	viper.SetEnvPrefix(EnvPrefix)
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
