playbook:
	go build && ./xp playbook --config  devopsxp.yaml

debug:
	go build && ./xp playbook -d --config  devopsxp.yaml

log:
	go build && ./xp playbook -d -l --config  devopsxp.yaml

cli:
	go build && ./xp cli shell 127.0.0.1-88 -u xp -a "hostname"

build:
	go build && ./xp -h

help:
	go run main.go -h

config:
	go build && ./xp config
