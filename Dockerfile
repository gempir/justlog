FROM quay.io/goswagger/swagger:latest as build-docs
WORKDIR /app
COPY . .
RUN make docs

FROM node:18-alpine as build-web
WORKDIR /web
COPY web .
COPY --from=build-docs /app/web/public/swagger.json /web/public
RUN yarn install --ignore-optional
RUN yarn build

FROM golang:alpine as build-app
WORKDIR /app
COPY . .
COPY --from=build-web /web/dist /app/web/dist
RUN go build -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build-app /app/app .
USER 1000:1000
CMD ["./app", "--config=/etc/justlog.json"]
EXPOSE 8025