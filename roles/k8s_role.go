package roles

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/devopsxp/xp/pkg/k8s"
	"github.com/devopsxp/xp/utils"
	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	// 初始化k8s role插件映射关系表
	addRoles(K8sType, reflect.TypeOf(K8sRole{}))
}

// 对标pod内containerd信息
type K8sBasic struct {
	command []string          // 执行命令
	env     map[string]string // 环境变量
	args    []string          // 命令参数
	name    string
	image   string
}

type GitRepo struct {
	url  string
	user string
	pwd  string
}

type K8sRole struct {
	RoleLC
	k8s       []K8sBasic // 容器pod组成
	workspace string     // 共享目录空间
	repo      GitRepo    // git代码
	name      string
	namespace string
}

func (k *K8sRole) Init(args *RoleArgs) error {
	err := k.Common(args)
	if err != nil {
		return err
	}

	// TODO: fixit
	k.name = "i-d-n"
	k.repo.url = args.reponame

	// 获取workspace
	if ws, ok := args.currentConfig["workspace"]; ok {
		k.workspace = ws.(string)
	} else {
		k.workspace = "/workspace"
	}

	// 获取namespace
	if ns, ok := args.currentConfig["namespace"]; ok {
		k.namespace = ns.(string)
	} else {
		k.namespace = "default"
	}

	// TODO: 解析k8s字段
	tmp := args.currentConfig["k8s"].([]interface{})
	for _, x := range tmp {
		k8sbasicData := K8sBasic{
			command: []string{},
			env:     map[string]string{},
			args:    []string{},
		}

		if n, ok := x.(map[interface{}]interface{})["name"]; ok {
			if n.(string) == "" {
				k8sbasicData.name = utils.GetRandomString(32)
			} else {
				k8sbasicData.name = n.(string)
			}
		}

		if im, ok := x.(map[interface{}]interface{})["image"]; ok {
			if im.(string) == "" {
				return errors.New("image is none")
			} else {
				k8sbasicData.image = im.(string)
			}
		}

		if sc, ok := x.(map[interface{}]interface{})["command"]; ok {
			for _, it := range sc.([]interface{}) {
				k8sbasicData.command = append(k8sbasicData.command, it.(string))
			}
		}

		if args, ok := x.(map[interface{}]interface{})["args"]; ok {
			for _, arg := range args.([]interface{}) {
				k8sbasicData.args = append(k8sbasicData.args, arg.(string))
			}
		}

		if e, ok := x.(map[interface{}]interface{})["env"]; ok {
			for k, v := range e.(map[interface{}]interface{}) {
				k8sbasicData.env[k.(string)] = v.(string)
			}
		}
		k.k8s = append(k.k8s, k8sbasicData)
	}

	return nil
}

func (k *K8sRole) Run() error {
	// 组装pod
	pod := &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:         k.name,
			GenerateName: k.name,
			Namespace:    k.namespace,
		},
		Spec: apiv1.PodSpec{
			Volumes: []apiv1.Volume{
				apiv1.Volume{
					Name: "workdir",
					VolumeSource: apiv1.VolumeSource{
						EmptyDir: &apiv1.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}

	containers := []apiv1.Container{}
	// TODO: env from k8srole => EnvVar
	for _, cc := range k.k8s {
		tmp := apiv1.Container{
			Name:    cc.name,
			Image:   cc.image,
			Command: cc.command,
			Args:    cc.args,
			VolumeMounts: []apiv1.VolumeMount{
				apiv1.VolumeMount{
					Name:      "workdir",
					MountPath: "/workspace",
				},
			},
		}
		containers = append(containers, tmp)
	}

	pod.Spec.Containers = containers

	log.Debugf("pod name: %s %v", k.name, pod)

	podinfo, err := k8s.CreatePod(pod)
	if err != nil {
		log.Error(err)
		return err
	}

	info, err := json.Marshal(podinfo)
	if err != nil {
		log.Error(err)
		return err
	}

	log.WithFields(log.Fields{
		"耗时": time.Now().Sub(k.starttime),
	}).Infof("Pod YAML: %s", info)
	return err
}

// 处理返回日志
func (k *K8sRole) After() {
	stoptime := time.Now()
	k.logs["耗时"] = fmt.Sprintf("%v", stoptime.Sub(k.starttime))
	k.msg.CallBack[fmt.Sprintf("%s-%s-%s", k.host, k.stage, k.name)] = k.logs
}
