// +build unit_tests

package service_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tamarakaufler/travel-article-headings/internal/service"

	"github.com/tamarakaufler/travel-article-headings/internal/photo"
)

func TestGetTopLocation(t *testing.T) {
	type args struct {
		articleLocation []photo.LocationM
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "location, weather and poi data fully provided",
			args: args{
				articleLocation: []photo.LocationM{
					{
						ArticleID: "AAA",
						PhotoID:   1,
						Location: photo.Location{
							Country: "USA",
							City:    "New York",
						},
					},
					{
						ArticleID: "AAA",
						PhotoID:   2,
						Location: photo.Location{
							Country: "Canada",
							City:    "Vancouver",
						},
					},
					{
						ArticleID: "AAA",
						PhotoID:   3,
						Location: photo.Location{
							Country: "Canada",
							City:    "Vancouver",
						},
					},
					{
						ArticleID: "AAA",
						PhotoID:   4,
						Location: photo.Location{
							Country: "USA",
							City:    "Sam Francisco",
						},
					},
					{
						ArticleID: "AAA",
						PhotoID:   5,
						Location: photo.Location{
							Country: "USA",
							City:    "New York",
						},
					},
					{
						ArticleID: "AAA",
						PhotoID:   5,
						Location: photo.Location{
							Country: "USA",
							City:    "San Francisco",
						},
					},
					{
						ArticleID: "AAA",
						PhotoID:   4,
						Location: photo.Location{
							Country: "USA",
							City:    "New York",
						},
					},
					{
						ArticleID: "AAA",
						PhotoID:   5,
						Location: photo.Location{
							Country: "USA",
							City:    "New York",
						},
					},
					{
						ArticleID: "AAA",
						PhotoID:   5,
						Location: photo.Location{
							Country: "USA",
							City:    "San Francisco",
						},
					},
				},
			},
			want:  "USA",
			want1: "New York",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := service.GetTopLocation(tt.args.articleLocation)
			require.NoError(t, err)

			if got != tt.want {
				t.Errorf("GetTopLocation() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetTopLocation() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
