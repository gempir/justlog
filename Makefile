full: web build

build: init_assets
	go build

run: build
	./justlog

run_web:
	cd web && yarn start

web: init_web
	cd web && yarn build

init_web:
	cd web && yarn install

init_assets:
	go run api/assets.go

# Docker stuff
container:
	docker build -t gempir/justlog .

release:
	docker push gempir/justlog

