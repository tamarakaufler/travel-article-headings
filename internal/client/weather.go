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

type Weather interface {
	EnhanceWithWeather(ctx context.Context, chans photo.Channel, pd photo.Data)
}

type weatherClient struct {
	client client
}

func NewWeatherClient(cfg conf.Setup) (weatherClient, error) {
	hereURL, err := url.Parse(cfg.WeatherURL)
	if err != nil {
		return weatherClient{}, err
	}

	return weatherClient{
		client: client{
			URL:    hereURL,
			APIKey: cfg.HereAPIKey,
			urlKey: "key",

			mu: &sync.Mutex{},
		},
	}, nil
}

var _ Weather = weatherClient{}

// EnhanceWithWeather fakes request for historical weather information.
func (wc weatherClient) EnhanceWithWeather(ctx context.Context, chans photo.Channel, pd photo.Data,
) {
	weatherL := []string{"rainy", "wet", "boiling hot", "sunny", "stormy", "drizzly", "hazy", "scorching",
		"hot", "unbearably hot", "miserably cold"}

	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(weatherL))

	w := photo.WeatherM{
		ArticleID: pd.ArticleID,
		PhotoID:   pd.ID,
		Weather:   weatherL[i],
	}

	ti, err := DateToSeason(pd.Date)
	if err == nil {
		w.TimeInfo = ti
	}

	chans.Weather <- w
}

// DateToSeason ...
func DateToSeason(d string) (photo.TimeInfo, error) {
	t, err := time.Parse(time.RFC3339, d)
	if err != nil {
		layout := "2006-01-02 15:04:05"
		t, err = time.Parse(layout, d)
		if err != nil {
			return photo.TimeInfo{}, err
		}
	}

	m := t.Month()
	wd := t.Weekday()

	mi := int(m)
	var s string
	switch {
	case mi >= 3 && mi < 6:
		s = "Spring"
	case mi >= 6 && mi < 9:
		s = "Summer"
	case mi >= 9 && mi < 12:
		s = "Autumn"
	default:
		s = "Winter"
	}

	return photo.TimeInfo{
		Weekday: wd.String(),
		Month:   m.String(),
		Season:  s,
	}, nil
}
