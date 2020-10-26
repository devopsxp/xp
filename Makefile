test:
	go build && ./xp test --config  devopsxp.yaml

debug:
	go build && ./xp test -d --config  devopsxp.yaml

log:
	go build && ./xp test -d -l --config  devopsxp.yaml

cli:
	go build && ./xp cli shell 127.0.0.1 192.168.50.1-255 -u lxp -a "hostname"

build:
	go build && ./xp -h

help:
	go run main.go -h

config:
	go build && ./xp config
