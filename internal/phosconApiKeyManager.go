package internal

import (
	"github.com/dgraph-io/badger/v3"
	"log"
)

const (
	DbApiKey = "apiKey"
)

type ApiKeyConfig struct {
	database   Operable
	httpClient *HttpClient
}

func NewApiKeyConfig(database Operable, httpClient *HttpClient) *ApiKeyConfig {
	return &ApiKeyConfig{
		database:   database,
		httpClient: httpClient,
	}
}

//TODO renvoyer les erreurs plutot que les logguer
func (config ApiKeyConfig) RegisterApiKey(gateway *Gateway) string {
	log.Println("Getting API Key from DB...")
	apiKey, err := config.database.Get(DbApiKey)
	var tmpApiKey string
	if err != nil {
		//TODO faire une erreur perso pour eviter l'import de badger
		if err == badger.ErrKeyNotFound {
			log.Printf("Key not found for %s", DbApiKey)
			log.Println("Trying to get the API Key...")
			jsonApiKey, err := config.httpClient.GetAndParseAPIKey(gateway)
			if err != nil {
				log.Fatal("Can't get API Key from Gateway", err)
			}
			tmpApiKey = jsonApiKey.Success.Username
			log.Println("Trying to insert the API Key in DB...")
			err = config.database.InsertOrUpdate(tmpApiKey, DbApiKey)
			if err != nil {
				log.Fatal("Can't insert API Key in DB", err)
			}
		} else {
			log.Fatal("Can't get and set apiKey", err)
		}
	} else {
		tmpApiKey = string(apiKey)
	}
	return tmpApiKey
}
