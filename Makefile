playbook:
	go build && ./xp playbook --config  devopsxp.yaml

debug:
	go build && ./xp playbook -d --config  devopsxp.yaml

log:
	go build && ./xp playbook -d -l --config  devopsxp.yaml

cli:
	go build && ./xp cli shell 127.0.0.1-88 -u lxp -a "hostname" -L console

shell:
	go build && ./xp cli shell 127.0.0.1 -a "for i in {1..100};do date;sleep 1;done" -u lxp -T

systemd:
	go build && ./xp cli systemd 127.0.0.1 -n docker -s status -u lxp

script:
	go build && ./xp cli script 127.0.0.1 -u lxp -a test.sh

copy:
	touch /tmp/123
	go build && ./xp cli copy 127.0.0.1 -u lxp -S /tmp/123 -D /tmp/333

fetch: clean
	touch /tmp/1abc
	mkdir /tmp/fetch
	go build && ./xp cli fetch 127.0.0.1 -u xp -S /tmp/1abc -D /tmp/fetch/2abc
	ls -lh /tmp/fetch

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

clean:
	rm -rf /tmp/fetch
	rm -f /tmp/1abc
