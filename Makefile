build:
	go build

full: build_web build_swagger build

full_run: full run

run:
	./justlog

run_web:
	cd web && yarn start

build_swagger:
	swag init

build_web: get_web
	cd web && npm run build
	go run api/assets.go

get_web:
	cd web && npm install

# Docker stuff
container:
	docker build -t gempir/justlog .

release:
	docker push gempir/justlog

