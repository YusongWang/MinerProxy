package fee

import (
	"net"
)

type FeeConn struct {
	DevConn net.Conn
	FeeConn net.Conn
}

// 保存临时份额 查找时间O(1)
type Fee struct {
	Dev map[string]bool
	Fee map[string]bool
}
