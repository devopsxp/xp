module github.com/devopsxp/xp

go 1.15

require (
	github.com/briandowns/spinner v1.11.1
	github.com/google/gops v0.3.13
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/lflxp/lflxp-kubectl v0.0.0-20200324075201-2c87c15ee42c
	github.com/mattbaird/jsonpatch v0.0.0-20200820163806-098863c1fc24
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/sftp v1.12.0
	github.com/satori/go.uuid v1.2.0
	github.com/shirou/gopsutil v3.20.10+incompatible // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/sys v0.0.0-20201126233918-771906719818 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.0.0-20201114085527-4a626d306b98
	k8s.io/apimachinery v0.0.0-20201118005411-2456ebdaba22
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20201114085527-4a626d306b98
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20201118005411-2456ebdaba22
)
