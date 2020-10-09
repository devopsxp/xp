package roles

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type RoleLC struct {
	name  string // 名称
	types string
	// 通用字段
	stage       string
	remote_user string                 // 执行用户
	vars        map[string]interface{} // 环境变量
	host        string                 // 执行的目标机
	starttime   time.Time              // 计算执行时间之开始时间
}

// 准备环节
func (r *RoleLC) Pre() {
	log.Debugf("Role module %s Pre running.", r.name)
	r.starttime = time.Now()
}

// 执行前
func (r *RoleLC) Before() {
	log.Debugf("Role module %s Before running.", r.name)
}

// 执行环节
func (r *RoleLC) Run() {
	log.Debugf("Role module %s Run running.", r.name)
}

// 执行后环节
func (r *RoleLC) After() {
	// log.Debugf("Role module %s After running.", r.name)
	log.WithFields(log.Fields{
		"Status": "After LifeCycle",
		"Role":   r.types,
		"Name":   r.name,
		"Stage":  r.stage,
		"Host":   r.host,
	}).Infof("执行完成 耗时：%v", time.Now().Sub(r.starttime))
}

// 执行判断IsHook
// default is false
func (r *RoleLC) IsHook() (string, string, bool) {
	return "", "", false
}

// 钩子函数，思考：是否和After以及output插件冲突
func (r *RoleLC) Hooks(types, target string, hook func(string, string) error) error {
	log.Debugf("Role module %s Hooks to %s:%s running.", r.name, types, target)
	err := hook(types, target)
	return err
}
