.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	go mod tidy
	env GOOS=linux go build -ldflags="-s -w" -o bin/mailer mailer/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/mailer mailer/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

test:
	go test ./mailer/

cover:
	go test -v -coverprofile=coverage.out ./mailer && go tool cover -html=coverage.out -o coverage.html