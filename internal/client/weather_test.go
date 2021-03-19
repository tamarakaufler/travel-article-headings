// +build unit_tests

package client_test

import (
	"reflect"
	"testing"

	"github.com/tamarakaufler/travel-article-headings/internal/client"
	"github.com/tamarakaufler/travel-article-headings/internal/photo"
)

func Test_dateToSeason(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name    string
		args    args
		want    photo.TimeInfo
		wantErr bool
	}{
		{
			name: "autumn test",
			args: args{
				d: "2019-10-29T11:11:59Z",
			},
			want: photo.TimeInfo{
				Weekday: "Tuesday",
				Month:   "October",
				Season:  "Autumn",
			},
			wantErr: false,
		},
		{
			name: "string test",
			args: args{
				d: "2019-03-29T11:11:59Z",
			},
			want: photo.TimeInfo{
				Weekday: "Friday",
				Month:   "March",
				Season:  "Spring",
			},
			wantErr: false,
		},
		{
			name: "summer test",
			args: args{
				d: "2019-07-29T11:11:59Z",
			},
			want: photo.TimeInfo{
				Weekday: "Monday",
				Month:   "July",
				Season:  "Summer",
			},
			wantErr: false,
		},
		{
			name: "winter test",
			args: args{
				d: "2019-12-29 11:11:59",
			},
			want: photo.TimeInfo{
				Weekday: "Sunday",
				Month:   "December",
				Season:  "Winter",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.DateToSeason(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("DateToSeason() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateToSeason() = %v, want %v", got, tt.want)
			}
		})
	}
}
