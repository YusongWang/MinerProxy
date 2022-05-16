package utils

import (
	"bytes"
	"encoding/gob"
	gomath "math"
	"math/big"
	"math/rand"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

var Ether = math.BigPow(10, 18)
var Shannon = math.BigPow(10, 9)

var pow256 = math.BigPow(2, 256)
var addressPattern = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
var zeroHash = regexp.MustCompile("^0?x?0+$")

func IsValidHexAddress(s string) bool {
	if IsZeroHash(s) || !addressPattern.MatchString(s) {
		return false
	}
	return true
}

func IsZeroHash(s string) bool {
	return zeroHash.MatchString(s)
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetTargetHex(diff int64) string {
	difficulty := big.NewInt(diff)
	diff1 := new(big.Int).Div(pow256, difficulty)
	return string(common.Bytes2Hex(diff1.Bytes()))
}

func TargetHexToDiff(targetHex string) *big.Int {
	targetBytes := common.FromHex(targetHex)
	return new(big.Int).Div(pow256, new(big.Int).SetBytes(targetBytes))
}

func ToHex(n int64) string {
	return "0x0" + strconv.FormatInt(n, 16)
}

func FormatReward(reward *big.Int) string {
	return reward.String()
}

func FormatRatReward(reward *big.Rat) string {
	wei := new(big.Rat).SetInt(Ether)
	reward = reward.Quo(reward, wei)
	return reward.FloatString(8)
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func MustParseDuration(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		panic("util: Can't parse duration `" + s + "`: " + err.Error())
	}
	return value
}

func String2Big(num string) *big.Int {
	n := new(big.Int)
	n.SetString(num, 0)
	return n
}

func BaseFeeToIndex(fee float64) uint64 {
	return uint64(int(gomath.Ceil(1000.0 / (fee * 10))))
}

func BaseOnIdxFee(idx uint64, fee float64) bool {
	return (idx % BaseFeeToIndex(fee)) == 0
}

func BaseOnRandFee(idx uint64, fee float64) bool {
	return rand.Intn(1000) <= int((fee+(fee*0.1))*10)
	//return rand.Intn(1000) <= int(fee*10)
}

func InterfaceToStrArray(list []interface{}) []string {
	job := make([]string, len(list))
	for i, arg := range list {
		job[i] = arg.(string)
	}

	return job
}

func DivTheDiff(newdiff *big.Int, olddiff *big.Int) *big.Int {
	if olddiff == new(big.Int).SetInt64(0) {
		return newdiff
	}
	return new(big.Int).Div(new(big.Int).Add(newdiff, olddiff), new(big.Int).SetInt64(2))
}

func IncreaseFDLimit() {
	var rlm syscall.Rlimit

	// Try to increase the soft limit
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlm)
	if rlm.Cur < 65535 && rlm.Cur < rlm.Max {
		rlm.Cur = rlm.Max
		syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlm)
	}

	// Try to increase the hard limit
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlm)
	if rlm.Cur < 65535 || rlm.Max < 65535 {
		rlm.Cur = 65535
		rlm.Max = 65535
		syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlm)
	}

	// checking
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlm)
	//Logger.Info("[OPTION] File descriptor limits")
	if rlm.Max < 5000 {
		Logger.Error("[OPTION] File descriptor hard limit is too small")
	}
	if rlm.Cur < 5000 {
		Logger.Error("[OPTION] File descriptor soft limit is too small")
	}
}

func HexRemovePrefix(hexStr string) string {
	// remove prefix "0x" or "0X"
	if len(hexStr) >= 2 && hexStr[0] == '0' && (hexStr[1] == 'x' || hexStr[1] == 'X') {
		hexStr = hexStr[2:]
	}
	return hexStr
}

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
