build:
	go build

full: build_web build_swagger build

full_run: full run

run: build
	./justlog

run_web:
	cd web && yarn start

build_swagger:
	swag init

build_web: get_web
	cd web && yarn build
	go run api/assets.go

get_web:
	cd web && yarn install

# Docker stuff
container:
	docker build -t gempir/justlog .

release:
	docker push gempir/justlog

