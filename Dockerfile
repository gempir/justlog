FROM golang:latest
WORKDIR /go/src/github.com/gempir/gempbotgo
RUN go get github.com/gempir/go-twitch-irc \
    && go get github.com/stretchr/testify/assert \
	&& go get github.com/labstack/echo \
	&& go get github.com/op/go-logging
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY configs ./configs
COPY --from=0 /go/src/github.com/gempir/gempbotgo/app .
CMD ["./app"]  
EXPOSE 8025
