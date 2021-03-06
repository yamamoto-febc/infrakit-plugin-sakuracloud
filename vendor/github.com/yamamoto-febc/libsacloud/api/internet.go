package api

import (
	"fmt"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"time"
)

type InternetAPI struct {
	*baseAPI
}

func NewInternetAPI(client *Client) *InternetAPI {
	return &InternetAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "internet"
			},
		},
	}
}

func (api *InternetAPI) UpdateBandWidth(id int64, bandWidth int) (*sacloud.Internet, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/bandwidth", api.getResourceURL(), id)
		body   = &sacloud.Request{}
	)
	body.Internet = &sacloud.Internet{BandWidthMbps: bandWidth}

	return api.request(func(res *sacloud.Response) error {
		return api.baseAPI.request(method, uri, body, res)
	})
}

func (api *InternetAPI) EnableIPv6(id int64) (*sacloud.IPv6Net, error) {
	var (
		method = "POST"
		uri    = fmt.Sprintf("%s/%d/ipv6net", api.getResourceURL(), id)
	)

	res := &sacloud.Response{}
	err := api.baseAPI.request(method, uri, nil, res)
	if err != nil {
		return nil, err
	}
	return res.IPv6Net, nil
}

func (api *InternetAPI) DisableIPv6(id int64, ipv6NetID int64) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%d/ipv6net/%d", api.getResourceURL(), id, ipv6NetID)
	)

	res := &sacloud.Response{}
	err := api.baseAPI.request(method, uri, nil, res)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (api *InternetAPI) SleepWhileCreating(internetID int64, timeout time.Duration) error {
	current := 0 * time.Second
	interval := 5 * time.Second

	var item *sacloud.Internet
	var err error
	//READ
	for item == nil && timeout > current {
		item, err = api.Read(internetID)

		if err != nil {
			time.Sleep(interval)
			current = current + interval
			err = nil
		}
	}

	if err != nil {
		return err
	}
	if current > timeout {
		return fmt.Errorf("Timeout: Can't read /internet/%d", internetID)
	}

	return nil

}

func (api *InternetAPI) Monitor(id int64, body *sacloud.ResourceMonitorRequest) (*sacloud.MonitorValues, error) {
	return api.baseAPI.monitor(id, body)
}
