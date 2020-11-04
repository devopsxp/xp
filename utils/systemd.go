package utils

/*
1. 匹配systemcd service or 其它
2. 统一封装服务的CRUD Function
3. switch type进行调用
*/
type Systemd interface {
	ListService() []map[string]string
	Search(string) map[string]string
	Start(string) error
	Restart(string) error
	Stop(string) error
	Status(string) error
}
