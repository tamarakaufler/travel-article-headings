package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/pkg/errors"
	"github.com/tamarakaufler/travel-article-headings/internal/client/response/here"
	conf "github.com/tamarakaufler/travel-article-headings/internal/configuration"
	"github.com/tamarakaufler/travel-article-headings/internal/photo"
)

type Addresses interface {
	EnhanceWithLocation(ctx context.Context, chans photo.Channel, pd photo.Data)
}

type addressesClient struct {
	client client
}

func NewAddressesClient(cfg conf.Setup) (addressesClient, error) {
	hereURL, err := url.Parse(cfg.HereURL)
	if err != nil {
		return addressesClient{}, err
	}

	return addressesClient{
		client: client{
			URL:    hereURL,
			APIKey: cfg.HereAPIKey,
			urlKey: "apiKey",

			mu: &sync.Mutex{},
		},
	}, nil
}

var _ Addresses = addressesClient{}

// EnhanceWithLocation ...
func (ac addressesClient) EnhanceWithLocation(ctx context.Context, chans photo.Channel, pd photo.Data,
) {
	ll := pd.LatLon
	at := latlonToAt(ll)
	q := map[string]string{
		"at":   at,
		"lang": "en-US",
	}

	//t1 := time.Now()
	b, err := ac.client.MakeGetRequest(ctx, q)
	if err != nil {
		chans.Error <- errors.Wrapf(err, "failure to retrieve location data for LatLon %+v", pd).Error()
		return
	}
	//t2 := time.Now()
	//log.Printf("\t\t--> call took %s", t2.Sub(t1).String())

	res := &here.ReverseGeocode{}
	err = json.Unmarshal(b, res)

	if err != nil {
		chans.Error <- errors.Wrapf(err, "failure to retrieve location data for LatLon %+v", pd).Error()
		return
	}
	if len(res.Items) == 0 {
		chans.Error <- fmt.Errorf("failure to retrieve location data for LatLon %+v", ll).Error()
		return
	}

	chans.Location <- photo.LocationM{
		ArticleID: pd.ArticleID,
		PhotoID:   pd.ID,
		Location: photo.Location{
			Country: res.Items[0].Address.CountryName,
			City:    res.Items[0].Address.City,
		},
	}
}

func latlonToAt(ll photo.LatLon) string {
	return fmt.Sprintf("%s,%s", ll.Latitude, ll.Longitude)
}
