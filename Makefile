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
	cd web && yarn install --ignore-optional

container:
	docker build -t gempir/justlog .

run_container:
	docker run -p 8025:8025 --restart=unless-stopped -v $(PWD)/config.json:/etc/justlog.json -v $(PWD)/logs:/logs gempir/justlog:latest

docs:
	swagger generate spec -m -o ./web/public/swagger.json -w api

