FROM golang:1.19-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

#RUN CGO_ENABLED=0 go test --tags=unit -v ./...

EXPOSE 2565

RUN go build -o ./out/go-app .

FROM alpine:3.16.2
COPY --from=build-base /app/out/go-app /app/go-app

CMD ["/app/go-app"]