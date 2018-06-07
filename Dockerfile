FROM golang:1.10 AS build_stage

WORKDIR /go/src/app

RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o watchkiller .

FROM scratch

COPY --from=build_stage /go/src/app/watchkiller /watchkiller
