.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/mailer mailer/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

test:
	go test ./mailer/

cover:
	go test -v -coverprofile=coverage.out ./mailer && go tool cover -html=coverage.out -o coverage.html