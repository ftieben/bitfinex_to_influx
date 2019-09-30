##################### gobuild #####################
FROM golang:alpine AS gobuild
RUN mkdir build
COPY . build/
WORKDIR build/
RUN apk add git
RUN go install
RUN CGO_ENABLED=0 GOOS=linux go build -o bitfinex bitfinex_to_influx.go

##################### run #####################
FROM alpine:latest
COPY --from=gobuild /go/build/bitfinex .
CMD ["/go/bitfinex"]