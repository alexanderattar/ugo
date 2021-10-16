FROM golang
ENV GO111MODULE=on
RUN mkdir -p $GOPATH/src/github.com/consensys
COPY . $GOPATH/src/github.com/consensys/ugo
WORKDIR $GOPATH/src/github.com/consensys/ugo
RUN go build -o ./cmd/api ./cmd/api
RUN go get github.com/rubenv/sql-migrate/...
CMD ["./cmd/api/api"]
