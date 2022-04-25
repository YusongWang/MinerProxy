package fee

import "io"

type FeeConn struct {
	DevConn *io.ReadWriteCloser
	FeeConn *io.ReadWriteCloser
}

// 保存临时份额 查找时间O(1)
type Fee struct {
	Dev map[string]bool
	Fee map[string]bool
}
