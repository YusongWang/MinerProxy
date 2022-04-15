package test

import (
	"testing"

	"github.com/buger/jsonparser"
)

var pushjob = `{"id":0,"version":"2.0","result":["0xf6067c77e474565f43b8c6b22afe7c12a178ec400e95974be5bad4f92657969d", "0xf87998fd030a4d04802b9f2ef04443bb4d2c105f13be5aa338ef49860c3c5425", "0x000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffff"]}`

var res = `{"id":0,"version":"2.0","result":["0xf6067c77e474565f43b8c6b22afe7c12a178ec400e95974be5bad4f92657969d", "0xf87998fd030a4d04802b9f2ef04443bb4d2c105f13be5aa338ef49860c3c5425", "0x000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffff"]}`

func TestJsonParse(t *testing.T) {
	data, type1, _, err := jsonparser.Get([]byte(pushjob), "result")
	if err != nil {
		t.Fatal(err.Error())
	}
	println(type1)
	println(string(data))
	t.Logf(string(data))
}
