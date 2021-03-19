package configuration_test

import (
	"os"
	"reflect"
	"testing"

	conf "github.com/tamarakaufler/travel-article-headings/internal/configuration"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		want    conf.Setup
		envs    map[string]string
		wantErr bool
	}{
		{
			name: "all conf set up",
			envs: map[string]string{
				"HERE_API_KEY":          "xxxxx",
				"GOOGLE_PLACES_API_KEY": "yyyyy",
				"WEATHER_API_KEY":       "zzzzz",
			},
			want: conf.Setup{
				Directory: "data4testing",

				HereURL:    "https://revgeocode.search.hereapi.com/v1/revgeocode",
				HereAPIKey: "xxxxx",

				GooglePlacesURL:    "https://maps.googleapis.com/maps/api/place/findplacefromtext/json",
				GooglePlacesAPIKey: "yyyyy",

				WeatherURL:    "https://weather.com/historical/json",
				WeatherAPIKey: "zzzzz",
			},
			wantErr: false,
		},
		{
			name:    "required env missing",
			envs:    map[string]string{},
			want:    conf.Setup{},
			wantErr: true,
		},
		{
			name: "custom conf set up",
			envs: map[string]string{
				"TRAVEL_ARTICLES_DIR":   "data",
				"HERE_URL":              "https://my.custom.url/json",
				"HERE_API_KEY":          "xxxxx",
				"GOOGLE_PLACES_API_KEY": "yyyyy",
				"WEATHER_API_KEY":       "zzzzz",
			},
			want: conf.Setup{
				Directory: "data",

				HereURL:    "https://my.custom.url/json",
				HereAPIKey: "xxxxx",

				GooglePlacesURL:    "https://maps.googleapis.com/maps/api/place/findplacefromtext/json",
				GooglePlacesAPIKey: "yyyyy",

				WeatherURL:    "https://weather.com/historical/json",
				WeatherAPIKey: "zzzzz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnvs(tt.envs)
			got, err := conf.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
			unSetEnvs(tt.envs)
		})
	}
}

func setEnvs(envsM map[string]string) {
	for k, v := range envsM {
		os.Setenv(k, v)
	}
}

func unSetEnvs(envsM map[string]string) {
	for k := range envsM {
		os.Unsetenv(k)
	}
}
