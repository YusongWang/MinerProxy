package utils

import (
	"testing"
)

func TestBaseOnIdxFee(t *testing.T) {
	ok := 0
	for i := 0; i < 1000; i++ {
		if BaseOnIdxFee(uint64(i), 0.3) {
			ok++
		}
	}

	if ok != 3 {
		t.Fatalf("Ok is Not 3 Ok:%v", ok)
	}

}
