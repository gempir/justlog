full: docs web build

build:
	go build

run: build
	./justlog

run_web:
	cd web && yarn start

web: init_web
	cd web && yarn build

init_web:
	cd web && yarn install

container:
	docker build -t gempir/justlog .

docs:
	swagger generate spec -m -o ./web/public/swagger.json -w api

# this is old fix later
#release:
#	docker push gempir/justlog

