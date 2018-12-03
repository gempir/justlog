build:
	go get ./... && env GOOS=linux GOARCH=amd64 go build

deploy: build
	ssh root@apollo.gempir.com systemctl stop justlog.service
	scp justlog root@apollo.gempir.com:/home/justlog/
	ssh root@apollo.gempir.com systemctl start justlog.service

provision: 
	ansible-playbook -i ansible/hosts ansible/playbook.yml --ask-vault-pass ${ARGS}