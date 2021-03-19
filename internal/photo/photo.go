package photo

import "sync"

type (
	// PhotoData ...
	Data struct {
		ArticleID string
		ID        int
		Date      string
		LatLon    LatLon
	}

	Location struct {
		Country string
		City    string
	}

	// LatLon ...
	LatLon struct {
		Latitude  string
		Longitude string
	}
)

// TimeInfo ...
type TimeInfo struct {
	Weekday string
	Month   string
	Season  string
}

// AdditionalInfo ...
type AdditionalInfo struct {
	City    string
	Country string

	TimeInfo

	Weather string

	// Example:
	//		["restaurant"]: 10,
	//		["bar"]: 12,
	//		["swimming_pool"]: 4 etc
	PlacesOfInterest map[string]int
}

type (
	LocationM struct {
		ArticleID string
		PhotoID   int
		Location  Location
	}
	WeatherM struct {
		ArticleID string
		PhotoID   int
		Weather   string
		TimeInfo  TimeInfo
	}
	PoiM struct {
		ArticleID string
		PhotoID   int
		POI       map[string]int
	}

	Channels map[string]Channel
	Channel  struct {
		Article  string
		Location chan LocationM
		Weather  chan WeatherM
		Poi      chan PoiM

		Cancel chan struct{}
		Error  chan string
	}
)

type (
	WgSyncs map[string]*WgSync
	WgSync  struct {
		Article string

		Location *sync.WaitGroup
		Weather  *sync.WaitGroup
		Poi      *sync.WaitGroup

		Heading *sync.WaitGroup
	}
)
