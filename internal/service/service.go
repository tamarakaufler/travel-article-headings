package service

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/tamarakaufler/travel-article-headings/internal/client"
	conf "github.com/tamarakaufler/travel-article-headings/internal/configuration"
	"github.com/tamarakaufler/travel-article-headings/internal/photo"
)

// ArticleHeadings ...
type ArticleHeadings []string

// Service interface prescribes methods that the service instance needs to implement.
type Service interface {
	GetArticles(ctx context.Context) ([]string, error)
	Run(ctx context.Context)
}

// ArticleService encapsulates service clients.
type ArticleService struct {
	Clients client.Clients
	Dir     string
}

// New is an ArticleService constructor.
func New(cfg conf.Setup, dir string) (ArticleService, error) {
	if dir == "" {
		dir = cfg.Directory
	}

	cs, err := client.BuildClients(cfg)
	if err != nil {
		return ArticleService{}, err
	}
	return ArticleService{
		Clients: cs,
		Dir:     dir,
	}, nil
}

var _ Service = ArticleService{}

//nolint:funlen
// Run processes the supplied csv files with photo date and lat/lon information
// and creates a list of heading suggestions for each file/article.
func (as ArticleService) Run(ctx context.Context) {
	albs, err := as.GetArticles(ctx)
	if err != nil {
		log.Fatalf(errors.Wrapf(err, "failure to get articles").Error())
	}
	albPaths := []string{}
	for _, alb := range albs {
		albPaths = append(albPaths, filepath.Join(as.Dir, alb))
	}
	chans, syncs := as.MakeChannelsAndSyncs(albPaths)

	doneCh := make(chan string)

	// listen and act
	go func(ctx context.Context, chans photo.Channels, done chan string) {
		for {
			for _, albP := range albPaths {
				select {
				case errM := <-chans[albP].Error:
					log.Printf("ERROR %s\n", errM)
				case <-chans[albP].Cancel:
					log.Print("unrecoverable error: cancelling request")
					os.Exit(1)
				case doneM := <-doneCh:
					log.Printf("%s\n", doneM)
					return
				default:
				}
			}
		}
	}(ctx, chans, doneCh)

	as.retrieveAdditionalData(ctx, albPaths, chans, syncs)
	ingestAdditionalInfoAndSuggestHeadings(ctx, albPaths, chans, syncs)

	// wait for additional photo info retrieval and processing.
	wgPF := &sync.WaitGroup{}

	// 	locations
	for _, albP := range albPaths {
		wgPF.Add(1)
		go func(wgPF *sync.WaitGroup, wgS *photo.WgSync, chanL chan photo.LocationM) {
			defer wgPF.Done()

			wgS.Location.Wait()
			close(chanL)
		}(wgPF, syncs[albP], chans[albP].Location)
	}

	// 	weather
	for _, albP := range albPaths {
		wgPF.Add(1)
		go func(wgPF *sync.WaitGroup, wgS *photo.WgSync, chanW chan photo.WeatherM) {
			defer wgPF.Done()

			wgS.Weather.Wait()
			close(chanW)
		}(wgPF, syncs[albP], chans[albP].Weather)
	}

	// 	poi
	for _, albP := range albPaths {
		wgPF.Add(1)
		go func(wgPF *sync.WaitGroup, wgS *photo.WgSync, chanP chan photo.PoiM) {
			defer wgPF.Done()

			wgS.Poi.Wait()
			close(chanP)
		}(wgPF, syncs[albP], chans[albP].Poi)
	}
	wgPF.Wait()

	// wait for heading creation and presentation to finish.
	wgTF := &sync.WaitGroup{}
	for _, albP := range albPaths {
		wgTF.Add(1)
		go func(wgTF *sync.WaitGroup, wgS *photo.WgSync) {
			defer wgTF.Done()

			wgS.Heading.Wait()
		}(wgTF, syncs[albP])
	}
	wgTF.Wait()

	doneCh <- "FINISHED 🎉"
}

// GetArticles provides a list of csv file paths.
func (as ArticleService) GetArticles(ctx context.Context) ([]string, error) {
	files, err := ioutil.ReadDir(as.Dir)
	if err != nil {
		return nil, err
	}

	csv := []string{}
	for _, f := range files {
		csv = append(csv, f.Name())
	}

	return csv, nil
}

func (as ArticleService) MakeChannelsAndSyncs(albs []string) (photo.Channels, photo.WgSyncs) {
	chans := photo.Channels{}
	syncs := photo.WgSyncs{}

	for _, alb := range albs {
		ch := photo.Channel{
			Article: alb,

			Location: make(chan photo.LocationM),
			Weather:  make(chan photo.WeatherM),
			Poi:      make(chan photo.PoiM),

			Cancel: make(chan struct{}),
			Error:  make(chan string),
		}
		chans[alb] = ch

		s := &photo.WgSync{
			Article: alb,

			Location: &sync.WaitGroup{},
			Weather:  &sync.WaitGroup{},
			Poi:      &sync.WaitGroup{},

			Heading: &sync.WaitGroup{},
		}
		syncs[alb] = s
	}

	return chans, syncs
}

func (as ArticleService) retrieveAdditionalData(ctx context.Context,
	albPaths []string, chans photo.Channels, syncs photo.WgSyncs) {
	for _, alb := range albPaths {
		photoL, err := ReadPhotoData(ctx, alb)
		if err != nil {
			log.Fatalf("failure to get photos %s", err)
		}

		// retrieval of location/weather/poi information.
		as.CollectAdditionalInfo(ctx, chans[alb], syncs[alb], photoL)
	}
}

// CollectAdditionalInfo retrieves additional photo info using 3rd party services.
func (as ArticleService) CollectAdditionalInfo(ctx context.Context,
	chans photo.Channel, wgS *photo.WgSync, photoL []photo.Data,
) {
	// retrieve location data for article photos.
	wgS.Location.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup, chans photo.Channel, photoL []photo.Data) {
		defer wg.Done()

		wgL := &sync.WaitGroup{}
		for _, pd := range photoL {
			wgL.Add(1)

			// Introducing backoff between goroutines to avoid HTTP 429/too many requests error/3rd party
			// rate limiting.
			// This is preferable to processing photos sequentially with a backoff between requests
			// as each request can take a different amount of time, so we can still take advantage of
			// concurrency.
			generateSleep(100, 190)

			go func(ctx context.Context, wgL *sync.WaitGroup, pd photo.Data) {
				defer wgL.Done()

				as.Clients.Addresses.EnhanceWithLocation(ctx, chans, pd)
			}(ctx, wgL, pd)
		}
		wgL.Wait()
	}(ctx, wgS.Location, chans, photoL)

	// retrieve weather data for article photos.
	wgS.Weather.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup, chans photo.Channel, photoL []photo.Data) {
		defer wg.Done()

		wgW := &sync.WaitGroup{}
		for _, pd := range photoL {
			wgW.Add(1)
			go func(ctx context.Context, wgW *sync.WaitGroup, pd photo.Data) {
				defer wgW.Done()

				as.Clients.Weather.EnhanceWithWeather(ctx, chans, pd)
			}(ctx, wgW, pd)
		}
		wgW.Wait()
	}(ctx, wgS.Weather, chans, photoL)

	// retrieve poi data for article photos.
	wgS.Poi.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup, chans photo.Channel, photoL []photo.Data) {
		defer wg.Done()

		wgP := &sync.WaitGroup{}
		for _, pd := range photoL {
			wgP.Add(1)
			go func(ctx context.Context, wgP *sync.WaitGroup, pd photo.Data) {
				defer wgP.Done()

				as.Clients.POI.EnhanceWithPlacesOfInterest(ctx, chans, pd)
			}(ctx, wgP, pd)
		}
		wgP.Wait()
	}(ctx, wgS.Poi, chans, photoL)
}

func ingestAdditionalInfoAndSuggestHeadings(ctx context.Context,
	albPaths []string, chans photo.Channels, syncs photo.WgSyncs,
) {
	for _, albP := range albPaths {
		wgT := syncs[albP].Heading
		chans := chans[albP]

		wgT.Add(1)
		go func(ctx context.Context, wgT *sync.WaitGroup, chans photo.Channel) {
			defer wgT.Done()

			articleLocationMap, articleWeatherMap, articlePoiMap := IngestAndProcess(ctx, chans.Article, chans)

			if len(articleLocationMap) <= 0 || len(articleWeatherMap) <= 0 || len(articlePoiMap) <= 0 {
				log.Print("Photo related data could not be retrieved.\nNo heading suggestions could be made.\n")
				log.Printf("locations: %d, weather data: %d, poi %d\n\n",
					len(articleLocationMap), len(articleWeatherMap), len(articlePoiMap))
				return
			}

			headings := CreateArticleHeadings(ctx,
				articleLocationMap, articleWeatherMap, articlePoiMap,
			)
			PresentSuggestedHeadings(chans.Article, headings)
		}(ctx, wgT, chans)
	}
}

// IngestAndProcess reads additional photo information from channels and
// processes it into input suited for article heading suggestions algorithm.
func IngestAndProcess(ctx context.Context, alb string, chans photo.Channel,
) ([]photo.LocationM, []photo.WeatherM, []photo.PoiM) {
	articleLocationMap := GatherLocationInfo(ctx, chans)

	// Uncomment below to see retrieved photo locations
	// log.Print("---------------------------------------\n")
	// log.Printf(">>> article %s\n\n", alb)
	// for _, v := range articleLocationMap {
	// 	log.Printf(">>> location %+v\n", v)
	// }
	// log.Print("---------------------------------------\n")

	articleWeatherMap := GatherWeatherInfo(ctx, chans)

	articlePoiMap := GatherPoiInfo(ctx, chans)

	return articleLocationMap, articleWeatherMap, articlePoiMap
}

func generateSleep(n, i int) {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(n) + i
	time.Sleep(time.Duration(r) * time.Millisecond) // to avoid HTTP 429, too many requests
}
