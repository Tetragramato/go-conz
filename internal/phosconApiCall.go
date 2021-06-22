package internal

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

const (
	SensorsUrl = "http://%s:%d/api/%s/sensors"
	ApiKeyUrl  = "http://%s:%d/api"
)

type httpClient struct {
	*resty.Client
}

func NewHttpClient() *httpClient {
	return &httpClient{Client: resty.New()}
}

func (client *httpClient) getRawAPIKey(gateway *Gateway) (*resty.Response, error) {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"devicetype":"go-conz"}`).
		Post(fmt.Sprintf(ApiKeyUrl, (*gateway)[0].Internalipaddress, (*gateway)[0].Internalport))

	if Config.TraceHttp {
		Trace(resp, err)
	}

	return resp, err
}

func (client *httpClient) GetGateway() (*Gateway, error) {
	resp, err := client.R().SetResult(&Gateway{}).
		Get(Config.PhosconUrl)

	if Config.TraceHttp {
		Trace(resp, err)
	}

	return resp.Result().(*Gateway), err
}

func (client *httpClient) getRawSensors(gateway *Gateway, apiKey string) (*resty.Response, error) {
	resp, err := client.R().
		Get(fmt.Sprintf(SensorsUrl, (*gateway)[0].Internalipaddress, (*gateway)[0].Internalport, apiKey))

	if Config.TraceHttp {
		Trace(resp, err)
	}

	return resp, err
}

func (client *httpClient) GetAndParseSensors(gatewayResp *Gateway) ([]*InputSensors, error) {
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
	listOfSensors, err := GetInputSensors(parsedJson)
	return listOfSensors, err
}
