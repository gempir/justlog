
default: dependencies gempbotgo 

.PHONY:
dependencies:
	go get github.com/gempir/go-twitch-irc
	go get github.com/stretchr/testify/assert
	go get github.com/labstack/echo

.PHONY: gempbotgo
gempbotgo:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gempbotgo .