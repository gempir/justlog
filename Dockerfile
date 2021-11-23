FROM node:16-alpine
WORKDIR /app
COPY . .
WORKDIR /app/web
RUN yarn install
RUN yarn build

FROM golang:latest
WORKDIR /app
COPY --from=0 /app .
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=1 /app/app .
CMD ["./app", "--config=/etc/justlog.json"]
EXPOSE 8025