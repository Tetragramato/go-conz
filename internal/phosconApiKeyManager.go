package internal

import (
	"encoding/json"
	"errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"time"
)

const (
	DbApiKey   = "apiKey"
	CountRetry = 10
)

var retryCounter int

type apiKeyConfig struct {
	database   Operable
	httpClient *httpClient
}

func NewApiKeyConfig(database Operable) *apiKeyConfig {
	httpClient := NewHttpClient()
	httpClient.
		SetRetryCount(CountRetry).
		SetRetryWaitTime(5 * time.Second).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				retryCounter++
				log.Printf("try (%d) a call to %s ...", retryCounter, r.Request.URL)
				return r.StatusCode() == http.StatusForbidden
			},
		)

	return &apiKeyConfig{
		database:   database,
		httpClient: httpClient,
	}
}

// RegisterApiKey TODO renvoyer les erreurs plutot que les logguer
func (config *apiKeyConfig) RegisterApiKey(gateway *Gateway) string {
	log.Println("Getting API Key from DB...")
	apiKey, err := config.database.Get(DbApiKey)
	var tmpApiKey string
	if err != nil {
		//TODO faire une erreur perso pour eviter l'import de badger
		if err == badger.ErrKeyNotFound {
			log.Printf("Key not found for %s", DbApiKey)
			log.Println("Trying to get the API Key...")
			jsonApiKey, err := config.getAndParseAPIKey(gateway)
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

func (config *apiKeyConfig) getAndParseAPIKey(gateway *Gateway) (*APIKey, error) {
	rawApiKey, err := config.httpClient.getRawAPIKey(gateway)
	if err != nil {
		return nil, err
	}

	if retryCounter > CountRetry {
		return nil, errors.New("fail to get APIKey : number of retries exceeded. Ensure you opened the Gateway to register a new application")
	}

	var parsedJson []interface{}
	err = json.Unmarshal(rawApiKey.Body(), &parsedJson)
	if err != nil {
		return nil, err
	}
	apiKey, err := GetApiKey(parsedJson)
	return apiKey, err
}
