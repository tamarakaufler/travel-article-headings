package client

import (
	"context"
	"math/rand"
	"net/url"
	"sync"
	"time"

	conf "github.com/tamarakaufler/travel-article-headings/internal/configuration"
	"github.com/tamarakaufler/travel-article-headings/internal/photo"
)

type Poi interface {
	EnhanceWithPlacesOfInterest(ctx context.Context, chans photo.Channel, pd photo.Data)
}

type poiClient struct {
	client client
}

func NewPoiClient(cfg conf.Setup) (poiClient, error) {
	hereURL, err := url.Parse(cfg.GooglePlacesURL)
	if err != nil {
		return poiClient{}, err
	}

	return poiClient{
		client: client{
			URL:    hereURL,
			APIKey: cfg.HereAPIKey,
			urlKey: "key",

			mu: &sync.Mutex{},
		},
	}, nil
}

var _ Poi = poiClient{}

// EnhanceWithPlacesOfInterest ...
func (pc poiClient) EnhanceWithPlacesOfInterest(ctx context.Context, chans photo.Channel, pd photo.Data,
) {
	placesL := []string{"Restaurants", "Casinos", "Museums", "Bars", "Swimming Pools", "Cafes", "Pubs",
		"Parks", "Theatres", "Cinemas", "Playgrounds", "Shopping Centres", "Zoos", "Botanical Gardens"}

	rand.Seed(time.Now().UnixNano())

	places := map[string]int{}
	for _, p := range placesL {
		c := rand.Intn(100)

		places[p] = c
	}

	chans.Poi <- photo.PoiM{
		ArticleID: pd.ArticleID,
		PhotoID:   pd.ID,
		POI:       places,
	}
}
