package internal

import (
	"encoding/json"
	"fmt"
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

func (config *apiKeyConfig) RegisterApiKey(gateway *Gateway) (string, error) {
	log.Println("Getting API Key from DB...")
	apiKey, err := config.database.Get(DbApiKey)
	var tmpApiKey string
	if err != nil {
		//TODO faire une erreur perso pour eviter l'import de badger
		if err == badger.ErrKeyNotFound {
			log.Printf("API Key not found for %s", DbApiKey)
			log.Println("Trying to get the API Key...")
			jsonApiKey, err := config.getAndParseAPIKey(gateway)
			if err != nil {
				return "", fmt.Errorf("can't get API Key from Gateway: %w", err)
			}
			tmpApiKey = jsonApiKey.Success.Username
			log.Println("Trying to insert the API Key in DB...")
			err = config.database.InsertOrUpdate(tmpApiKey, DbApiKey)
			if err != nil {
				return "", fmt.Errorf("can't insert API Key in DB: %w", err)
			}
		} else {
			return "", fmt.Errorf("can't get and set API Key: %w", err)
		}
	} else {
		tmpApiKey = string(apiKey)
	}
	return tmpApiKey, nil
}

func (config *apiKeyConfig) getAndParseAPIKey(gateway *Gateway) (*APIKey, error) {
	rawApiKey, err := config.httpClient.getRawAPIKey(gateway)
	if err != nil {
		return nil, err
	}

	if retryCounter > CountRetry {
		return nil, fmt.Errorf("number of retries exceeded (%v). Ensure you opened the Gateway to register a new application", CountRetry)
	}

	var parsedJson []interface{}
	err = json.Unmarshal(rawApiKey.Body(), &parsedJson)
	if err != nil {
		return nil, err
	}
	apiKey, err := GetApiKey(parsedJson)
	return apiKey, err
}
