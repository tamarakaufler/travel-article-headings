// +build service_tests

package service_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tamarakaufler/travel-article-headings/internal/client"
	"github.com/tamarakaufler/travel-article-headings/internal/photo"
	"github.com/tamarakaufler/travel-article-headings/internal/service"
)

type addrClientM struct {
}
type weathClientM struct {
}
type poiClientM struct {
}

func Test_CollectAdditionalInfo(t *testing.T) {
	// prepare test data.
	ctx := context.Background()
	photoL := []photo.Data{
		{
			ArticleID: "article1",
			ID:        1,
			Date:      "2019-11-25T22:37:44Z",
			LatLon: photo.LatLon{
				Latitude:  "36.111111",
				Longitude: "16.11111",
			},
		},
		{
			ArticleID: "article1",
			ID:        2,
			Date:      "2019-10-31T11:26:42Z",
			LatLon: photo.LatLon{
				Latitude:  "36.222222",
				Longitude: "16.22222",
			},
		},
		{
			ArticleID: "article1",
			ID:        3,
			Date:      "2019-11-25T22:37:44Z",
			LatLon: photo.LatLon{
				Latitude:  "36.333333",
				Longitude: "16.33333",
			},
		},
	}
	as, ch, s := setup("article1")

	wgF := &sync.WaitGroup{}
	wgF.Add(1)
	go func(wgF *sync.WaitGroup, s *photo.WgSync) {
		defer wgF.Done()

		s.Location.Wait()
	}(wgF, s)

	wgF.Add(1)
	go func(wgF *sync.WaitGroup, s *photo.WgSync) {
		defer wgF.Done()

		s.Weather.Wait()
	}(wgF, s)

	wgF.Add(1)
	go func(wgF *sync.WaitGroup, s *photo.WgSync) {
		defer wgF.Done()

		s.Poi.Wait()
	}(wgF, s)

	// run the method.
	as.CollectAdditionalInfo(ctx, ch, s, photoL)

	// drain the channels and verify.
	timer := time.After(3 * time.Second)

	msgL := len(photoL)
	loopC := 0
	locCount := 0
	weathCount := 0
	poiCount := 0
LOOP:
	for {
		select {
		case <-timer:
			t.Error("Test_CollectAdditionalInfo timed out")
			break LOOP
		case <-ch.Location:
			locCount++
		case <-ch.Weather:
			weathCount++
		case <-ch.Poi:
			poiCount++
		}
		loopC++

		if loopC >= msgL*3 {
			break LOOP
		}
	}
	wgF.Wait()

	require.Equal(t, msgL, locCount)
	require.Equal(t, msgL, weathCount)
	require.Equal(t, msgL, poiCount)
}

func setup(alb string) (service.ArticleService, photo.Channel, *photo.WgSync) {
	ac := addrClientM{}
	wc := weathClientM{}
	pc := poiClientM{}

	as := service.ArticleService{
		Clients: client.Clients{
			Addresses: ac,
			Weather:   wc,
			POI:       pc,
		},
		Dir: "somedir",
	}

	ch := createArticleChannel(alb)
	s := createSync(alb)

	return as, ch, s
}

func createArticleChannel(alb string) photo.Channel {
	return photo.Channel{
		Article: alb,

		Location: make(chan photo.LocationM),
		Weather:  make(chan photo.WeatherM),
		Poi:      make(chan photo.PoiM),

		Cancel: make(chan struct{}),
		Error:  make(chan string),
	}
}

func createSync(alb string) *photo.WgSync {
	return &photo.WgSync{
		Article: alb,

		Location: &sync.WaitGroup{},
		Weather:  &sync.WaitGroup{},
		Poi:      &sync.WaitGroup{},

		Heading: &sync.WaitGroup{},
	}
}

func (ac addrClientM) EnhanceWithLocation(ctx context.Context, chans photo.Channel, pd photo.Data) {
	chans.Location <- photo.LocationM{
		ArticleID: pd.ArticleID,
		PhotoID:   pd.ID,
		Location: photo.Location{
			Country: "Czechia",
			City:    "Prague",
		},
	}
}
func (wc weathClientM) EnhanceWithWeather(ctx context.Context, chans photo.Channel, pd photo.Data,
) {
	chans.Weather <- photo.WeatherM{
		ArticleID: pd.ArticleID,
		PhotoID:   pd.ID,
		Weather:   "sunny",
	}
}

func (pc poiClientM) EnhanceWithPlacesOfInterest(ctx context.Context, chans photo.Channel, pd photo.Data,
) {
	places := map[string]int{
		"Cinemas":     5,
		"Restaurants": 10,
		"Cafes":       15,
	}
	chans.Poi <- photo.PoiM{
		ArticleID: pd.ArticleID,
		PhotoID:   pd.ID,
		POI:       places,
	}
}
