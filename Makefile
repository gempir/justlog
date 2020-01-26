build: get
	go build

full: build_web build_swagger build

get:
	go get ./... 

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

