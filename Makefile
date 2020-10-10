test:
	go build && ./xp test --config  devopsxp.yaml

debug:
	go build && ./xp test -d --config  devopsxp.yaml

log:
	go build && ./xp test -d -l --config  devopsxp.yaml

build:
	go build && ./xp -h

config:
	go build && ./xp config