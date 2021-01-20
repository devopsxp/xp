// 主机磁盘配置获取数据结构
package setup

import "github.com/shirou/gopsutil/disk"

type Disk struct {
	disk.PartitionStat
}
