package switchbot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// getDevicesResp represents a response from GET /devices API.
// https://github.com/OpenWonderLabs/SwitchBotAPI/blob/main/README.md#devices
type getDevicesResp struct {
	StatusCode int                `json:"statusCode"`
	Message    string             `json:"message"`
	Body       getDevicesRespBody `json:"body"`
}

type getDevicesRespBody struct {
	DeviceList         []device         `json:"deviceList"`
	InfraredRemoteList []infraredRemote `json:"infraredRemoteList"`
}

type device struct {
	ID                 string `json:"deviceId"`
	Name               string `json:"deviceName"`
	Type               string `json:"deviceType"`
	EnableCloudService bool   `json:"enableCloudService"`
	HubDeviceId        string `json:"hubDeviceId"`

	// TODO: Add missing fields defined in the API document.
}

type infraredRemote struct {
	ID          string `json:"deviceId"`
	Name        string `json:"deviceName"`
	Type        string `json:"remoteType"`
	HubDeviceId string `json:"hubDeviceId"`
}

func (c *Client) GetDevices(ctx context.Context) (*getDevicesResp, error) {
	resp, err := c.do(ctx, http.MethodGet, "/devices")
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	var response getDevicesResp
	if err := json.NewDecoder(resp).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// getDeviceStatusResp represents a response from GET /devices/{deviceId}/status API.
// https://github.com/OpenWonderLabs/SwitchBotAPI/blob/main/README.md#get-device-status
type getDeviceStatusResp struct {
	StatusCode int          `json:"statusCode"`
	Message    string       `json:"message"`
	Body       deviceStatus `json:"body"`
}

type deviceStatus struct {
	ID          string  `json:"deviceId"`
	Type        string  `json:"deviceType"`
	Power       string  `json:"power"`
	HubDeviceId string  `json:"hub_device_id"`
	Temperature float64 `json:"temperature"`
	Huidity     int     `json:"humidity"`

	// TODO: Add missing fields defined in the API document.
}

func (c *Client) GetDeviceStatus(ctx context.Context, deviceID string) (*getDeviceStatusResp, error) {
	resp, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/devices/%s/status", deviceID))
	if err != nil {
		return nil, fmt.Errorf("failed to do the request: %w", err)
	}
	defer resp.Close()

	var response getDeviceStatusResp
	if err := json.NewDecoder(resp).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
