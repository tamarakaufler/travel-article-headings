
deps:
	@go mod download
	@go mod tidy

lint:
	@golangci-lint -v run

test:
	go test -v -count=1 --race -tags unit_tests -cover ./...
	go test -v -count=1 -tags service_tests -cover ./...

test_client:
	go test -v -count=1 --race -tags client_tests -cover ./...

build:
	CGO_ENABLED=0 go build -o ./cmd/bin/travel-article-headings ./cmd/travel-article-headings/main.go

bin:
	./cmd/bin/travel-article-headings

run:
	go run ./cmd/travel-article-headings/main.go

# Running through docker image

## Local docker image

docker-build:
	docker build -t travel-article-headings:v1.0.0 .

docker-run-default:
	docker run --rm -v ${PWD}/data4testing:/data4testing --env HERE_API_KEY=xxxx travel-article-headings:v1.0.0

docker-kill:
	docker ps | grep "travel-article-headings:v1.0.0" | awk '{print $1}'  | xargs docker kill
	
## Remote (prebuilt) docker image

docker-run-custom:
	docker run --rm -v ${PWD}/data:/data -v ${PWD}/data4testing:/data4testing --env HERE_API_KEY=xxxx --env TRAVEL_ARTICLES_DIR=data travel-article-headings:v1.0.0

all: deps lint test build bin

.PHONY: deps lint test test_client build bin run