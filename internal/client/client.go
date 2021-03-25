package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	conf "github.com/tamarakaufler/travel-article-headings/internal/configuration"
)

type client struct {
	URL    *url.URL
	APIKey string
	urlKey string

	mu *sync.Mutex
}

// Clients collect all the required clients.
type Clients struct {
	Addresses Addresses
	Weather   Weather
	POI       Poi
}

// BuildClients ...
func BuildClients(cfg conf.Setup) (Clients, error) {
	addressesClient, err := NewAddressesClient(cfg)
	if err != nil {
		return Clients{}, err
	}

	weatherClient, err := NewWeatherClient(cfg)
	if err != nil {
		return Clients{}, err
	}

	poiClient, err := NewPoiClient(cfg)
	if err != nil {
		return Clients{}, err
	}

	return Clients{
		Addresses: addressesClient,
		Weather:   weatherClient,
		POI:       poiClient,
	}, nil
}

// MakeGetRequest will be used for all 3rd party requests.
func (c client) MakeGetRequest(ctx context.Context, query map[string]string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	u := c.URL
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	q.Set(c.urlKey, c.APIKey)
	u.RawQuery = q.Encode()

	req := &http.Request{
		Method: "GET",
		URL:    u,
	}
	req.Header = headers()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusTooManyRequests {
			return nil, fmt.Errorf("client (%s) responds with too many requests error", c.URL.String())
		}
		return nil, fmt.Errorf("failure to retrieve photo location: HTTP status = %d", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func headers() map[string][]string {
	return map[string][]string{
		"Accept": {"application/json"},
		// This header results in "invalid character '\x1f' looking for beginning of value"
		// (https://github.com/influxdata/influxdb-go/issues/12)
		//	"Accept-Encoding": {"gzip", "deflate", "br"}:
		//					causes an error "invalid character '\x1f' looking for beginning of value".
		//					The reason is that, if using this header, the 3rd party sends compressed
		//					content, that would need to be uncompressed before unmarshalling.
		"Content-Type":    {"application/json; charset=UTF-8"},
		"Connection":      {"keep-alive"},
		"Cache-Control":   {"max-age=0"},
		"User-Agent":      {"travel-article-headings"},
		"Accept-Language": {"en-GB", "en-US"},
	}
}
