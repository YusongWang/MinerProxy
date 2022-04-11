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

	ok = 0
	for i := 0; i < 1000; i++ {
		if BaseOnIdxFee(uint64(i), 1) {
			ok++
		}
	}

	if ok != 10 {
		t.Fatalf("Ok is Not 3 Ok:%v", ok)
	}

	ok = 0
	for i := 0; i < 1000; i++ {
		if BaseOnIdxFee(uint64(i), 10) {
			ok++
		}
	}

	if ok != 100 {
		t.Fatalf("Ok is Not 3 Ok:%v", ok)
	}

	ok = 0
	for i := 0; i < 1000; i++ {
		if BaseOnIdxFee(uint64(i), 50) {
			ok++
		}
	}

	if ok != 500 {
		t.Fatalf("Ok is Not 3 Ok:%v", ok)
	}

	ok = 0
	for i := 0; i < 1000; i++ {
		if BaseOnIdxFee(uint64(i), 1000) {
			ok++
		}
	}

	if ok != 1000 {
		t.Fatalf("Ok is Not 3 Ok:%v", ok)
	}
}
