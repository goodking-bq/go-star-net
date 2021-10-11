package conn

import "sync"

// 网络
// 网络下面有很多节点
type NetWorker interface {
	Name() string      // 网络名唯一
	NetMask() string   // netmask
	Nodes() sync.Map   // 节点信息
	CurrentNode() Node // 当前节点
}
