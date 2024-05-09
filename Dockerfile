FROM golang:1.18-alpine

RUN apk add g++ && apk add make


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

EXPOSE 8080
ENV PORT=8080

CMD ["make", "run"]
