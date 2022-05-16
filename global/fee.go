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

// TODO 保存当前JOBS
// type DevJob struct {
// 	Job [][]string
// }

type Job struct {
	Target string
	JobId  string
	Diff   string
	Job    []byte
}
