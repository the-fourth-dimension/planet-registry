dev:
	export APP_ENV=DEV;\
	go run cmd/main.go

build:
	go build -o planet-registry-bin cmd/main.go

run: build
	export APP_ENV=PRODUCTION;\
	./planet-registry-bin

test:
	export APP_ENV=TEST;\
	go test -v ./tests/...