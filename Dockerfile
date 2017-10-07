FROM golang:latest
WORKDIR /go/src/github.com/gempir/gempbotgo
COPY . .
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/gempir/gempbotgo/app .
CMD ["./app"]  
EXPOSE 8025
