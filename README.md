# planet registry

planet registry for the hyperspace

## Running locally

```sh
# setup env vars
cp .env.example .env

# run the planet-registry dev server
go run cmd/main.go

# build the planet-registry server
go build -o planet-registry-bin cmd/main.go

# run the build
./planet-registry-bin

# Running tests
go test -v ./tests/...
```
