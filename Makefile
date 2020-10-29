playbook:
	go build && ./xp playbook --config  devopsxp.yaml

debug:
	go build && ./xp playbook -d --config  devopsxp.yaml

log:
	go build && ./xp playbook -d -l --config  devopsxp.yaml

cli:
	go build && ./xp cli shell 127.0.0.1-88 -u lxp -a "hostname"

copy:
	touch /tmp/123
	go build && ./xp cli copy 127.0.0.1 -u lxp -S /tmp/123 -D /tmp/333

template:
	echo "hello {{.data}}" > /tmp/tmp.j2
	go build && ./xp cli template 127.0.0.1 -u lxp -S template.service.j2  -D /tmp/docker.service
	rm -f /tmp/tmp.j2

build:
	go build && ./xp -h

help:
	go run main.go -h

config:
	go build && ./xp config
