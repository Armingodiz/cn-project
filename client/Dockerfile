FROM golang:1.16-alpine AS go-build

ENV GO111MODULE=on

WORKDIR /app

COPY ./client/go.mod .
COPY ./client/go.sum .

RUN go mod download

COPY ./client .






RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o client .

FROM alpine:3.9.5 as dns
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=go-build /app .
ENTRYPOINT ["./client"]