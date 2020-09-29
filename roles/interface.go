package roles

// 执行过程生命周期接口
type RoleLifeCycle interface {
	// 环境准备
	Pre()

	// 执行前
	Before()

	// 执行中
	// 返回是否执行信号
	Run() error

	// 执行后
	After()

	// 是否执行hook
	IsHook() (string, string, bool)

	// 钩子函数
	Hooks(string, string, func(string, string) error) error
}
