test:
	go build && ./xp test --config  devopsxp.backup.yaml

debug:
	go build && ./xp test -d --config  devopsxp.backup.yaml

log:
	go build && ./xp test -d -l --config  devopsxp.backup.yaml

build:
	go build && ./xp -h

config:
	go build && ./xp config