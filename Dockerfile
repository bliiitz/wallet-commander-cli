FROM golang:1.17 as common
# All these steps will be cached
RUN mkdir /wallet-commander
WORKDIR /wallet-commander
COPY go.mod . 
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# COPY the source code as the last step
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -extldflags "-static"' -o /go/bin/wallet-commander-cli
RUN chmod +x /go/bin/wallet-commander-cli
ENTRYPOINT ["/go/bin/wallet-commander-cli"]