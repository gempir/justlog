build: get
	swag init
	go build

get:
	go get ./... 

container:
	docker build -t gempir/justlog .

release:
	docker push gempir/justlog

