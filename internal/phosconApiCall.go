package internal

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

const (
	SensorsUrl = "http://%s:%d/api/%s/sensors"
	ApiKeyUrl  = "http://%s:%d/api"
)

func GetAPIKey(client *resty.Client, gateway Gateway) (*APIKey, error) {

	resp, err := client.R().SetResult(&APIKey{}).
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"devicetype":"go-conz"}`).
		Post(fmt.Sprintf(ApiKeyUrl, gateway[0].Internalipaddress, gateway[0].Internalport))
	if err != nil {
		return nil, err
	}

	if TraceHttp{
		Trace(resp, err)
	}

	return resp.Result().(*APIKey), nil
}

func GetGateway(client *resty.Client) (*Gateway, error) {
	resp, err := client.R().SetResult(&Gateway{}).
		EnableTrace().
		Get(PhosconUrl)
	if err != nil {
		return nil, err
	}

	if TraceHttp{
		Trace(resp, err)
	}

	return resp.Result().(*Gateway), nil

}

func GetSensors(client *resty.Client, gateway *Gateway, apiKey string) (*resty.Response, error) {
	resp, err := client.R().
		EnableTrace().
		Get(fmt.Sprintf(SensorsUrl, (*gateway)[0].Internalipaddress, (*gateway)[0].Internalport, apiKey))
	if err != nil {
		return nil, err
	}

	if TraceHttp{
		Trace(resp, err)
	}

	return resp, nil
}
