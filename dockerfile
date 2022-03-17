FROM golang:latest as builder

# Go directories
ENV GOPATH=/go
RUN mkdir -p /go/src && mkdir -p /go/bin && mkdir -p /go/pkg && mkdir -p /go/src/build
ENV PATH=$GOPATH/bin:$PATH

# Install dep
RUN go get -u github.com/golang/dep/...

# Copy sources
COPY /server /go/src/build
WORKDIR /go/src/build

# Go dep
RUN dep ensure

# Go build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o baywatch .

# Output image
FROM alpine:3.15
COPY --from=builder /go/src/build /app/
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/* && update-ca-certificates
WORKDIR /app
CMD ["./baywatch"]