// +build client_tests

//nolint:testpackage
package client

import (
	"context"
	"encoding/json"
	"net/url"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	here "github.com/tamarakaufler/travel-article-headings/internal/client/response/here"
	conf "github.com/tamarakaufler/travel-article-headings/internal/configuration"
)

func Test_client_MakeGetRequest(t *testing.T) {
	type args struct {
		ctx   context.Context
		query map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    *here.ReverseGeocode
		wantErr bool
	}{
		{
			name: "address retrieval - happy path",
			args: args{
				ctx: context.Background(),
				query: map[string]string{
					"at":   "40.628075,14.375383",
					"lang": "en-US",
				},
			},
			want:    unmarshalExpectedResponse(t, "response/here/hereRevgeocodeResponse.json"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs, err := conf.Load()
			require.NoError(t, err)

			u, err := url.Parse(cs.HereURL)
			require.NoError(t, err)

			c := client{
				URL:    u,
				APIKey: cs.HereAPIKey,
				urlKey: "apiKey",
				mu:     &sync.Mutex{},
			}
			got, err := c.MakeGetRequest(tt.args.ctx, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.MakeGetRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)

			res, err := unmarshalReverseGeocodeResponse(got)
			require.NoError(t, err)

			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("client.MakeGetRequest() = %v, want %v", res, tt.want)
			}
		})
	}
}

func unmarshalReverseGeocodeResponse(b []byte) (*here.ReverseGeocode, error) {
	res := &here.ReverseGeocode{}
	err := json.Unmarshal(b, res)
	return res, err
}

func unmarshalExpectedResponse(t *testing.T, path string) *here.ReverseGeocode {
	b, err := os.ReadFile(path)
	require.NoError(t, err)

	data := &here.ReverseGeocode{}
	err = json.Unmarshal(b, data)
	require.NoError(t, err)

	return data
}
