FROM golang:1.13-stretch as scratch
WORKDIR /src

COPY . .

RUN  GOARCH=amd64 GOOS=linux go build ./cmd/main.go

FROM scratch
COPY --from=scratch /src/main /bin/main

ENTRYPOINT ["/bin/main"]
