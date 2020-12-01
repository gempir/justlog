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

container:
	docker build -t gempir/justlog .

docs:
	swagger generate spec -o ./api/swagger.json -w api

# this is old fix later
#release:
#	docker push gempir/justlog

