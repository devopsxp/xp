package utils

import (
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

type dockerCli struct {
	image   string
	args    []string // 启动参数
	command string   // 最终执行命令
	cmd     string   // docker命令
}

func NewDockerCli(args []string, image, cmd string) *dockerCli {
	return &dockerCli{
		args:  args,
		image: image,
		cmd:   cmd,
	}
}

// 查询docker命令是否存在
func (d *dockerCli) CheckPath() error {
	if path, err := exec.LookPath("docker"); err != nil {
		return err
	} else {
		d.command = fmt.Sprintf("%s run -it --rm ", path)
	}

	return nil
}

// 通过struct数据进行args聚合
func (d *dockerCli) AddArgs() *dockerCli {
	if len(d.args) <= 0 {
		return d
	}

	if d.command != "" {
		for _, arg := range d.args {
			d.command += arg + " "
		}
	}
	return d
}

// 通过arg实时添加args
// arg eg: -v /tmp/123:/data
func (d *dockerCli) AddArg(arg string) *dockerCli {
	if arg != "" {
		d.command = fmt.Sprintf("%s %s ", d.command, arg)
	}

	return d
}

func (d *dockerCli) Run() error {
	if err := d.CheckPath(); err != nil {
		return err
	}

	d.AddArgs()

	log.Warnln(fmt.Sprintf("%s %s %s", d.command, d.image, d.cmd))

	rs, err := ExecCommandString(fmt.Sprintf("%s %s %s", d.command, d.image, d.cmd))
	if err != nil {
		return err
	}

	log.Info(rs)
	return nil
}
