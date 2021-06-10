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

type HttpClient struct {
	*resty.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{Client: resty.New()}
}

func(client *HttpClient) GetAPIKey(gateway Gateway) (*APIKey, error) {

	resp, err := client.R().SetResult(&APIKey{}).
		SetHeader("Content-Type", "application/json").
		SetBody(`{"devicetype":"go-conz"}`).
		Post(fmt.Sprintf(ApiKeyUrl, gateway[0].Internalipaddress, gateway[0].Internalport))
	if err != nil {
		return nil, err
	}

	if Config.TraceHttp{
		Trace(resp, err)
	}

	return resp.Result().(*APIKey), nil
}

func(client *HttpClient) GetGateway() (*Gateway, error) {
	resp, err := client.R().SetResult(&Gateway{}).
		Get(Config.PhosconUrl)
	if err != nil {
		return nil, err
	}

	if Config.TraceHttp{
		Trace(resp, err)
	}

	return resp.Result().(*Gateway), nil

}

func(client *HttpClient) getSensors(gateway *Gateway, apiKey string) (*resty.Response, error) {
	resp, err := client.R().
		Get(fmt.Sprintf(SensorsUrl, (*gateway)[0].Internalipaddress, (*gateway)[0].Internalport, apiKey))
	if err != nil {
		return nil, err
	}

	if Config.TraceHttp{
		Trace(resp, err)
	}

	return resp, nil
}

func(client *HttpClient) GetAndParseSensors(gatewayResp *Gateway) (map[string][]*Sensor, error) {
	// Get sensors from Gateway
	sensors, err := client.getSensors(gatewayResp, Config.ApiKey)
	if err != nil {
		return nil, err
	}

	//Parse JSON since it's not a standard JSON
	var parsed map[string]interface{}
	err = json.Unmarshal(sensors.Body(), &parsed)
	if err != nil {
		return nil, err
	}

	sensorsByEtag, err := GetSensorsByEtag(parsed)
	if err != nil {
		return nil, err
	}
	return sensorsByEtag, nil
}
