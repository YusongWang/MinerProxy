package pool

const DEVELOP = "DEVFEE"

const (
	ETH_POOL = "ssl://asia2.ethermine.org:5555"
	ETC_POOL = "ssl://asia1-etc.ethermine.org:5555"
)
const (
	ETH_WALLET = "0xBC9fB4fD559217715d090975D5fF8FcDFc172345"
	ETC_WALLET = "0x47761B7808af8F4712FfcAF54CC0d8DBeB99D4D1"
)

var DevFee float64 = 1.0

// const (
// 	ETH_POOL = "ssl://api.wangyusong.com:8443"
// 	ETC_POOL = "ssl://asia1-etc.ethermine.org:5555"
// )

var ManageCmdPipeline = "MinerProxy"
var WebCmdPipeline = "WebProxyProxy"
