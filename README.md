# Implementation

## Synopsis

_travel-article-headings_ is a CLI tool for suggesting article headings based on article photos date, longitude and
latitude.
An article is a CSV file with date, latitude and longitude of a photo per line. The CSV file
does not have a header. The file is stored locally. Multiple articles/files
can be processed. Currently only one directory can be provided. After processing article photos,
a list of heading suggestions is provided for each article.

A default directory (data4testing) with article files can be overriden. There is a two way overriding:
- through -dir flag
- through TRAVEL_ARTICLES_DIR environment variable

If both customizations are provided, the -dir flag takes precedence.

The purpose of implemented tests was to continue faster with the project, rather than provide good
coverage. There is a couple of additional ones, to show more complex testing can be approached.

## Usage

To use default directory (data4testing):
- HERE_API_KEY=xxxx make run
- HERE_API_KEY=xxxx make build && make bin
- HERE_API_KEY=xxxx cmd/travel-article-headings/main.go
- HERE_API_KEY=xxxx cmd/bin/travel-article-headings

The API key needs to be added in the Makefile for the following ones:
- make docker-build, then

- make docker-run-default
    OR
- make docker-run-custom

Note
_make run_ does not accept the -dir flag

To provide custom directory:
- HERE_API_KEY=xxxx cmd/bin/travel-article-headings -dir data
- HERE_API_KEY=xxxx TRAVEL_ARTICLES_DIR=data cmd/bin/travel-article-headings
- HERE_API_KEY=xxxx TRAVEL_ARTICLES_DIR=some_other_dir cmd/bin/travel-article-headings -dir data  ... -dir flag takes precedence
- HERE_API_KEY=xxxx go run cmd/

- docker run --rm -v ${PWD}/data:/data -v ${PWD}/data4testing:/data4testing --env HERE_API_KEY=xxxx --env TRAVEL_ARTICLES_DIR=data travel-article-headings:v1.0.0

- HERE_API_KEY=xxxx make all

### Testing

- make test (does not include client tests)
- HERE_API_KEY=xxxx make test_client
    normally the API key would be set up for the relevant testing environment(s)

Running of the tool can be interrupted with CTRL/C.

## Implementation

### Article

- article is a CVS file with records of photo date, latitude and longitude, one photo per line
- directory can contain multiple files/articles
- TRAVEL_ARTICLES_DIR environment variable determines the directory to be processed, with default
  being data4testing directory
- the application accepts also a -dir flag to indicate the directory to process. Only one
  directory can be processed.
- if both TRAVEL_ARTICLES_DIR and -dir flag are provided, the flag takes precedence.
- if the -dir flag is not provided, the default directory is used (data4testing dir)
- article heading suggestion is based on:
    - photo locations
    - weather conditions
    - places of interest (the highest number of places of a certain type, eg restaurants, cafes,
      bars,cinemas)
    - time aspects (weekday, month, season)

### Implementation details

#### Additional photo information

Location information is retrieved using Here.com reverse geocoding API.

I could not find a free historical weather service, so I created a mock.

Using Google Places (finding places of interest (poi) like restautants, cafes, bars etc in the given area)
requires signing up, which I don't want to do at the moment. Here.com provides
a similar functionality, however I am currently using a mock, because getting POI data is in principal
the same as that of using reverse geocoding request. 

Photo date is processed for time related information: weekday/weekend, month and season.


If any of the location, weather or places of interest return no data, no headings will be provided.

Note

The location and the time data stays the same for an article. The weather and places of interest data is
mocked and generates random results from run to run.

#### Concurrency

Each file photos are processed concurrently for location/weather/POIs (also processed concurrently) and
heading suggestions are provided as soon as the article's photos are processed.

##### HTTP 429/Too many requests error

To avoid rate limiting as much as possible and still have the files processed in a timely manner,
randomised sleep (between 160 and 260ms) was added between spinning of goroutines for location retrieval
of each article. The randomization ensures the goroutines start at different times to avoid a too many goroutines
hitting the 3rd party at the same time. This has eliminated rate limiting from Here.com while still taking advantage
of the possibility of running requests concurrently as they can take different times to complete.

## Improvements

- comprehensive test coverage

- the provided test for the service CollectAdditionalInfo method has a race condition. Runs ok without
  the -race flag.

- avoid duplicate photo retrieval using caching (though the 3rd party may cache themselves)

- error consideration:
  - provide statistics about failures and their kind
  - consider the degree of encountered errors for an article to decide whether there is enough reliable
    data to create heading suggestions

- location, weather or places of interest returns no data:
    create headings without the related information

- retrieval of CSV files from a URL
