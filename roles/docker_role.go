package roles

import (
	"fmt"
	"reflect"
	"time"

	"github.com/devopsxp/xp/utils"
	log "github.com/sirupsen/logrus"
)

func init() {
	// 初始化docker role插件映射关系表
	addRoles(DockerType, reflect.TypeOf(DockerRole{}))
}

/* @Comment: 原型YAML功能点

# 定义 job
Build:
  stage: build
  image: node:8.15.1-jessie
  script:
    - pwd
    - whoami

    - npm version
    - npm install -g @vue/cli --registry=https://registry.npm.taobao.org
    - npm install --registry=https://registry.npm.taobao.org
    - npm run build --target=prod

    - tar -zcvf dist.tar.gz dist/
  artifacts:
    name: "${CI_JOB_STAGE}_${CI_COMMIT_REF_NAME}"
    expire_in: 3 day
    paths:
      - dist.tar.gz
  only:
    - master
  tags:
    # - 10-128-6-109
    - k8s_public_centos7.2.1511
    # - k8s_public_centos_runner_12_4
  retry: 2
*/

type DockerRole struct {
	RoleLC
	script []string // 执行脚本命令
	image  string   // 执行镜像
}

// 准备数据
// @Param stage 阶段标记
// @Param user 远端执行用户
// @Param host 目标主机
// @Param vars 动态参数
// @Param configs 执行模块内容
// @Param msg 消息结构体
func (r *DockerRole) Init(args *RoleArgs) error {
	err := r.Common(args)
	if err != nil {
		return err
	}

	// 获取镜像
	r.image = args.currentConfig["image"].(string)

	// 获取script迭代
	if sc, ok := args.currentConfig["script"]; ok {
		for _, it := range sc.([]interface{}) {
			r.script = append(r.script, it.(string))
		}
	}

	return nil
}

func (r *DockerRole) Run() error {
	// Docker RemoteAPI Dedefault Port
	port := "9999"
	env := []string{}

	for k, v := range r.vars {
		if reflect.TypeOf(v).String() == "string" {
			env = append(env, fmt.Sprintf("%s=%s", k, v.(string)))
		}
	}

	data := map[string]interface{}{
		"Hostname":     "",
		"Domainname":   "",
		"User":         "",
		"AttachStdin":  false,
		"AttachStdout": true,
		"AttachStderr": true,
		"Tty":          false,
		"OpenStdin":    false,
		"StdinOnce":    false,
		"Env":          env,
		"Cmd":          r.script,
		"Entrypoint":   "",
		"Image":        r.image,
		"Volumes": map[string]interface{}{
			"/tmp": map[string]string{},
		},
		"WorkingDir":      "/",
		"NetworkDisabled": false,
		"ExposedPorts": map[string]interface{}{
			"22/tcp": map[string]string{},
		},
		"StopSignal": "SIGTERM",
		"HostConfig": map[string]interface{}{
			"Binds":              []string{"/tmp:/tmp"},
			"Tmpfs":              map[string]string{"/run": "rw,noexec,nosuid,size=65536k"},
			"Links":              []string{}, // "redis3:redis"
			"Memory":             0,          // 8MB
			"MemorySwap":         0,
			"MemoryReservation":  0,
			"KernelMemory":       0,
			"CpuShares":          512,
			"CpuPeriod":          100000,
			"CpuQuota":           50000,
			"CpusetCpus":         "",
			"CpusetMems":         "",
			"IOMaximumBandwidth": 0,
			"IOMaximumIOps":      0,
			"MemorySwappiness":   60,
			"OomKillDisable":     false,
			"OomScoreAdj":        500,
			"PidMode":            "",
			"PidsLimit":          -1,
			"PortBindings":       map[string]interface{}{"22/tcp": []map[string]string{map[string]string{"HostPort": "11022"}}},
			"PublishAllPorts":    false,
			"Privileged":         false,
			"ReadonlyRootfs":     false,
			"Dns":                []string{"8.8.8.8"},
			"DnsOptions":         []string{},
			"DnsSearch":          []string{},
			"ExtraHosts":         []string{},
			"VolumesFrom":        []string{}, // ["parent", "other:ro"],
			"CapAdd":             []string{"NET_ADMIN"},
			"CapDrop":            []string{"MKNOD"},
			"RestartPolicy":      map[string]interface{}{"Name": "", "MaximumRetryCount": 0},
			"NetworkMode":        "bridge",
			"Devices":            []string{},
			"Sysctls":            map[string]string{"net.ipv4.ip_forward": "1"},
			"Ulimits":            []map[string]string{},
			"LogConfig":          map[string]interface{}{"Type": "json-file", "Config": map[string]string{}},
			"SecurityOpt":        []string{},
			"CgroupParent":       "",
			"VolumeDriver":       "",
			"ShmSize":            67108864,
		},
	}

	cli := utils.NewDockerCLI(r.host, port, "")
	rs, err := cli.CreateContainer(data)
	if err != nil {
		log.WithFields(log.Fields{
			"Host":   r.host,
			"Name":   r.name,
			"Script": len(r.script),
			"Stage":  r.stage,
			"User":   r.remote_user,
			"耗时":     time.Now().Sub(r.starttime),
		}).Errorln(err.Error())
		r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = err.Error()
		return err
	}

	log.WithFields(log.Fields{
		"Host":   r.host,
		"Name":   r.name,
		"Script": len(r.script),
		"Stage":  r.stage,
		"User":   r.remote_user,
		"耗时":     time.Now().Sub(r.starttime),
	}).Info(string(rs))
	r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = string(rs)

	return nil
}

// 处理返回日志
func (r *DockerRole) After() {
	stoptime := time.Now()
	r.logs["耗时"] = fmt.Sprintf("%v", stoptime.Sub(r.starttime))
	r.msg.CallBack[fmt.Sprintf("%s-%s-%s", r.host, r.stage, r.name)] = r.logs
}
