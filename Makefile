full: web swag build

build:
	go build

run: build
	./justlog

run_web:
	cd web && yarn start

swag: init_swag init_assets

web: init_web
	cd web && yarn build

init_web:
	cd web && yarn install

init_swag:
	swag init -g api/server.go --output web/public
	rm web/public/docs.go
	rm web/public/swagger.yaml

init_assets:
	go run api/assets.go

# Docker stuff
container:
	docker build -t gempir/justlog .

release:
	docker push gempir/justlog

