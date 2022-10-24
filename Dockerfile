#build stage
FROM golang:alpine 
RUN apk add --no-cache git make 

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go mod download


