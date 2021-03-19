package configuration

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

// Setup holds:
//		3rd Party API URL and key information.
type Setup struct {
	Directory  string `env:"TRAVEL_ARTICLES_DIR" envDefault:"data4testing"`
	HereURL    string `env:"HERE_URL" envDefault:"https://revgeocode.search.hereapi.com/v1/revgeocode"` // prox=x.x,y.y&mode=retrieveAddresses
	HereAPIKey string `env:"HERE_API_KEY,required"`

	// ==> will simulate
	GooglePlacesURL    string `env:"GOOGLE_PLACES_URL" envDefault:"https://maps.googleapis.com/maps/api/place/findplacefromtext/json"` // input=address
	GooglePlacesAPIKey string `env:"GOOGLE_PLACES_API_KEY" envDefault:"yyyyy"`

	// ==> will simulate
	WeatherURL    string `env:"WEATHER_URL" envDefault:"https://weather.com/historical/json"`
	WeatherAPIKey string `env:"WEATHER_API_KEY" envDefault:"zzzzz"`
}

// Load customizes configuration based on env variables.
func Load() (Setup, error) {
	cfg := &Setup{}

	for _, c := range []interface{}{
		cfg,
	} {
		if err := env.Parse(c); err != nil {
			return Setup{}, errors.Wrapf(err, "failed to load configuration")
		}

	}
	return *cfg, nil
}
