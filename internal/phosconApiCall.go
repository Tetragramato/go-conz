package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"time"
)

const (
	SensorsUrl = "http://%s:%d/api/%s/sensors"
	ApiKeyUrl  = "http://%s:%d/api"
	CountRetry = 10
)

type HttpClient struct {
	*resty.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{Client: resty.New()}
}

func (client *HttpClient) getRawAPIKey(gateway *Gateway) (*resty.Response, error) {

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"devicetype":"go-conz"}`).
		Post(fmt.Sprintf(ApiKeyUrl, (*gateway)[0].Internalipaddress, (*gateway)[0].Internalport))

	if Config.TraceHttp {
		Trace(resp, err)
	}

	return resp, err
}

func (client HttpClient) GetAndParseAPIKey(gateway *Gateway) (*APIKey, error) {
	var retryCounter int
	client.
		SetRetryCount(CountRetry).
		SetRetryWaitTime(5 * time.Second).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				retryCounter++
				log.Printf("try (%d) a call to %s ...", retryCounter, r.Request.URL)
				return r.StatusCode() == http.StatusForbidden
			},
		)
	rawApiKey, err := client.getRawAPIKey(gateway)
	if err != nil {
		return nil, err
	}

	if retryCounter > CountRetry {
		return nil, errors.New("fail to get APIKey : number of retries exceeded. Ensure you opened the Gateway to register a new application")
	}
	//TODO voir si c'est utile
	client.SetRetryCount(0)

	var parsedJson []interface{}
	err = json.Unmarshal(rawApiKey.Body(), &parsedJson)
	if err != nil {
		return nil, err
	}
	apiKey, err := GetApiKey(parsedJson)
	return apiKey, err
}

func (client *HttpClient) GetGateway() (*Gateway, error) {
	resp, err := client.R().SetResult(&Gateway{}).
		Get(Config.PhosconUrl)

	if Config.TraceHttp {
		Trace(resp, err)
	}

	return resp.Result().(*Gateway), err
}

func (client *HttpClient) getRawSensors(gateway *Gateway, apiKey string) (*resty.Response, error) {
	resp, err := client.R().
		Get(fmt.Sprintf(SensorsUrl, (*gateway)[0].Internalipaddress, (*gateway)[0].Internalport, apiKey))

	if Config.TraceHttp {
		Trace(resp, err)
	}

	return resp, err
}

func (client *HttpClient) GetAndParseSensors(gatewayResp *Gateway) ([]*SensorsList, error) {
	// Get sensors from Gateway
	rawSensors, err := client.getRawSensors(gatewayResp, Config.ApiKey)
	if err != nil {
		return nil, err
	}

	//Parse JSON since it's not a standard JSON
	var parsedJson map[string]interface{}
	err = json.Unmarshal(rawSensors.Body(), &parsedJson)
	if err != nil {
		return nil, err
	}

	listOfSensorsList, err := GetListOfSensorsList(parsedJson)
	return listOfSensorsList, err
}
