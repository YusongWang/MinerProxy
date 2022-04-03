package fee

import "sync"

// 保存临时份额 查找时间O(1)
type Fee struct {
	Dev map[string]bool
	Fee map[string]bool
	sync.RWMutex
}
