package service

import (
	"context"
	"encoding/csv"
	"log"
	"os"

	"github.com/tamarakaufler/travel-article-headings/internal/photo"
)

// Article ...
type Article struct {
	Name      string
	PhotoData []photo.Data
	Headings  []string
}

// ReadPhotoData ...
func ReadPhotoData(ctx context.Context, fp string) ([]photo.Data, error) {
	log.Printf("Processing article: %s\n", fp)

	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	photoD := []photo.Data{}
	for i, r := range records {
		pd := photo.Data{
			ArticleID: fp,
			ID:        i + 1,
			Date:      r[0],
			LatLon: photo.LatLon{
				Latitude:  r[1],
				Longitude: r[2],
			},
		}
		photoD = append(photoD, pd)
	}

	return photoD, nil
}
