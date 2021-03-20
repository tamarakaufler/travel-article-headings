package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/tamarakaufler/travel-article-headings/internal/photo"
)

type SortPositionList []SortPosition
type SortPosition struct {
	name  string
	count int
}

// GatherLocationInfo gathers photo info for each article.
func GatherLocationInfo(ctx context.Context, chans photo.Channel,
) []photo.LocationM {
	articleLocationData := []photo.LocationM{}

	//nolint:gosimple
	for {
		select {
		case m, ok := <-chans.Location:
			if !ok {
				return articleLocationData
			}
			articleLocationData = append(articleLocationData, m)
		}
	}
}

// GatherWeatherInfo gathers photo info for each article.
//func GatherWeatherInfo(ctx context.Context, weatherCh chan photo.WeatherM, doneWeatherCh chan struct{},
func GatherWeatherInfo(ctx context.Context, chans photo.Channel) []photo.WeatherM {
	articleWeatherData := []photo.WeatherM{}

	//nolint:gosimple
	for {
		select {
		case m, ok := <-chans.Weather:
			if !ok {
				return articleWeatherData
			}
			articleWeatherData = append(articleWeatherData, m)
		}
	}
}

// GatherPoiInfo gathers photo info for each article.
func GatherPoiInfo(ctx context.Context, chans photo.Channel,
) []photo.PoiM {
	articlePOIData := []photo.PoiM{}

	//nolint:gosimple
	for {
		select {
		case m, ok := <-chans.Poi:
			if !ok {
				return articlePOIData
			}
			articlePOIData = append(articlePOIData, m)
		}
	}
}

//nolint:funlen
// CreateArticleHeadings - this is where the fun magic happens.
// The function processes one article.
func CreateArticleHeadings(ctx context.Context,
	articleLocationData []photo.LocationM, articleWeatherData []photo.WeatherM, articlePOIData []photo.PoiM,
) ArticleHeadings {
	headingStart1 := []string{"Enjoy your break", "Having great time",
		"Experience of a lifetime", "Have a holiday of a lifetime", "Wonderful break"}
	headingStart2 := []string{"Enjoy happy days", "Have hilarious time"}
	headingMiddle := []string{"with friends", "with family", "on your own"}
	headingForPlaces := []string{"full of", "bursting with", "brimming with", "packed with"}
	adjectives := []string{"Hilarious", "Beautiful", "Brilliant", "Family fun", "Glorious"}

	country, city, errLoc := GetTopLocation(articleLocationData)
	if errLoc != nil {
		city = "city"
		country = "country"
	}
	weather := GetTopWeather(articleWeatherData)
	weekday, month, season := GetTopTimeInfo(articleWeatherData)
	if weekday == "Saturday" || weekday == "Sunday" {
		weekday = "Weekend"
	}

	poi := GetTopPlaceOfInterest(articlePOIData)

	headings := []string{}
	for i := 0; i < 1; i++ {
		rand.Seed(time.Now().UnixNano())
		jj := rand.Intn(len(headingStart1))
		kk := rand.Intn(len(headingMiddle))

		s := headingStart1[jj]
		m := headingMiddle[kk]
		heading := strings.Join([]string{s, m, "in", weather, country}, " ")
		headings = append(headings, heading)
	}

	for i := 0; i < 1; i++ {
		rand.Seed(time.Now().UnixNano())
		jj := rand.Intn(len(adjectives))
		kk := rand.Intn(len(headingMiddle))

		s := adjectives[jj]
		m := headingMiddle[kk]
		heading := strings.Join([]string{s, month, m, "in", weather, country}, " ")
		headings = append(headings, heading)
	}

	for i := 0; i < 1; i++ {
		rand.Seed(time.Now().UnixNano())
		jj := rand.Intn(len(adjectives))
		mm := rand.Intn(len(headingForPlaces))
		p := headingForPlaces[mm]

		s := adjectives[jj]
		heading := strings.Join([]string{s, weekday, "enjoying", city, "of", p, poi}, " ")
		headings = append(headings, heading)
	}

	for i := 0; i < 1; i++ {
		rand.Seed(time.Now().UnixNano())
		jj := rand.Intn(len(adjectives))
		mm := rand.Intn(len(headingForPlaces))
		p := headingForPlaces[mm]

		s := adjectives[jj]
		heading := strings.Join([]string{s, season, "break in", country, p, poi}, " ")
		headings = append(headings, heading)
	}

	for i := 0; i < 1; i++ {
		rand.Seed(time.Now().UnixNano())
		jj := rand.Intn(len(adjectives))
		mm := rand.Intn(len(headingForPlaces))
		p := headingForPlaces[mm]

		s := adjectives[jj]
		heading := strings.Join([]string{s, weekday, "stay in", city, p, poi}, " ")
		headings = append(headings, heading)
	}

	for i := 0; i < 1; i++ {
		rand.Seed(time.Now().UnixNano())
		jj := rand.Intn(len(headingStart2))
		kk := rand.Intn(len(headingMiddle))

		s := headingStart2[jj]
		m := headingMiddle[kk]
		heading := strings.Join([]string{s, m, "in", weather, city}, " ")
		headings = append(headings, heading)
	}

	for i := 0; i < 1; i++ {
		rand.Seed(time.Now().UnixNano())
		jj := rand.Intn(len(headingStart1))
		kk := rand.Intn(len(headingMiddle))
		mm := rand.Intn(len(headingForPlaces))

		s := headingStart1[jj]
		m := headingMiddle[kk]
		p := headingForPlaces[mm]
		heading := strings.Join([]string{s, m, "in", country, p, poi}, " ")
		headings = append(headings, heading)
	}

	for i := 0; i < 1; i++ {
		rand.Seed(time.Now().UnixNano())
		jj := rand.Intn(len(headingStart1))
		kk := rand.Intn(len(headingForPlaces))

		s := headingStart1[jj]
		p := headingForPlaces[kk]
		heading := strings.Join([]string{s, "in", city, p, poi}, " ")
		headings = append(headings, heading)
	}

	return headings
}

// GetTopLocation ...
func GetTopLocation(articleLocation []photo.LocationM) (string, string, error) {
	if len(articleLocation) <= 0 {
		return "", "", nil
	}

	// determine the number of occurrencies.
	countryM := map[string]int{}
	cityM := map[string]int{}
	for _, l := range articleLocation {
		countryM[l.Location.Country] = countryM[l.Location.Country] + 1
		cityM[l.Location.City] = cityM[l.Location.City] + 1
	}

	// find the most frequent occurrence of a country.
	topLocationL := SortPositionList{}
	for k, v := range countryM {
		topLocationL = append(topLocationL, SortPosition{
			name:  k,
			count: v,
		})
	}
	sort.Sort(sort.Reverse((topLocationL)))
	topCountry := topLocationL[0].name

	// find the most frequent occurrence of a city.
	topLocationL = SortPositionList{}
	for k, v := range cityM {
		topLocationL = append(topLocationL, SortPosition{
			name:  k,
			count: v,
		})
	}
	sort.Sort(sort.Reverse((topLocationL)))
	topCity := topLocationL[0].name

	return topCountry, topCity, nil
}

// GetTopWeather ...
func GetTopWeather(weatherData []photo.WeatherM) string {
	// determine the number of occurrencies.
	weather := map[string]int{}
	for _, w := range weatherData {
		weather[w.Weather] = weather[w.Weather] + 1
	}

	// find the most frequent occurrence of weather.
	topWeatherL := SortPositionList{}
	for k, v := range weather {
		topWeatherL = append(topWeatherL, SortPosition{
			name:  k,
			count: v,
		})
	}
	sort.Sort(sort.Reverse((topWeatherL)))
	topWeather := topWeatherL[0].name

	return topWeather
}

// GetTopTimeInfo ...
func GetTopTimeInfo(weatherData []photo.WeatherM) (string, string, string) {
	// determine the number of occurrencies.
	weekday := map[string]int{}
	month := map[string]int{}
	season := map[string]int{}
	for _, w := range weatherData {
		ti := w.TimeInfo
		weekday[ti.Weekday] = weekday[ti.Weekday] + 1
		month[ti.Month] = month[ti.Month] + 1
		season[ti.Season] = season[ti.Season] + 1
	}

	// find the most frequent occurrence of weekday.
	topWeekdayL := SortPositionList{}
	for k, v := range weekday {
		topWeekdayL = append(topWeekdayL, SortPosition{
			name:  k,
			count: v,
		})
	}
	sort.Sort(sort.Reverse((topWeekdayL)))
	topWeekday := topWeekdayL[0].name

	// find the most frequent occurrence of month.
	topMonthL := SortPositionList{}
	for k, v := range month {
		topMonthL = append(topMonthL, SortPosition{
			name:  k,
			count: v,
		})
	}
	sort.Sort(sort.Reverse((topMonthL)))
	topMonth := topMonthL[0].name

	// find the most frequent occurrence of season.
	topSeasonL := SortPositionList{}
	for k, v := range season {
		topSeasonL = append(topSeasonL, SortPosition{
			name:  k,
			count: v,
		})
	}
	sort.Sort(sort.Reverse((topSeasonL)))
	topSeason := topSeasonL[0].name

	return topWeekday, topMonth, topSeason
}

func GetTopPlaceOfInterest(poiData []photo.PoiM) string {
	// determine the number of occurrencies.
	poi := map[string]int{}
	for _, p := range poiData {
		pp := p.POI
		for k, v := range pp {
			poi[k] = poi[k] + v
		}
	}

	// find the most frequent occurrence of poi.
	topPOIL := SortPositionList{}
	for k, v := range poi {
		topPOIL = append(topPOIL, SortPosition{
			name:  k,
			count: v,
		})
	}
	sort.Sort(sort.Reverse((topPOIL)))
	topPOI := topPOIL[0].name

	return topPOI
}

// PresentSuggestedHeadings ...
func PresentSuggestedHeadings(alb string, headings ArticleHeadings) {
	fmt.Printf("---------------------------------------\n")
	log.Printf("%s\n\n", alb)
	for _, t := range headings {
		fmt.Printf("\t%s\n", t)
	}
	fmt.Printf("---------------------------------------\n")
}

// Custom sorting to determine the highest ranking of Photo attributes.
func (spl SortPositionList) Len() int           { return len(spl) }
func (spl SortPositionList) Less(i, j int) bool { return spl[i].count < spl[j].count }
func (spl SortPositionList) Swap(i, j int)      { spl[i], spl[j] = spl[j], spl[i] }
