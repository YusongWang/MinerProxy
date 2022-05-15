package global

import (
	"io"
	"sync"
)

type FeeConn struct {
	DevConn *io.ReadWriteCloser
	FeeConn *io.ReadWriteCloser
}

type FeeResult struct{}

// 保存临时份额 查找时间O(1)
type Fee struct {
	Dev sync.Map
	Fee sync.Map
}

// 保存当前JOB
type Job struct {
	Job [][]string
}
