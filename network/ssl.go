package network

import (
	"crypto/tls"
	"miner_proxy/utils"
	"net"
)

const inner_cert = `-----BEGIN CERTIFICATE-----
MIIFWzCCA0OgAwIBAgIUfJgXhRMNlglWAOYEQLuO2VX6lBMwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMjA0MDIwNjA5MzhaFw0yMzA0
MDIwNjA5MzhaMFoxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxEzARBgNVBAMMCmZvb2Jh
ci5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC97F2tSGuZfd+C
sT2fckKJ4YCNpKf5P2fcWlqalU0uAqJiH55cFWssDcuFc6AHv9cflcK8iUqMXs1c
JzI25c6aTTbkdooVcUkAeHoRx9iruO/q7qhHmKLgnCQX/FfhNBad2Csuit/T4mk5
SCa07KZmChRjUeXZJQ/OUghWPxTlr3SnsbOJqrQTnWZUjEh8IIEQ+v1YhB8gYOS3
d594Qn69qbMNZHSY0xAn/nLoCDJ1zgJq/5UuwszV556U+A/aG81ypGaq27PCXujM
jlbuAs/c6BzOSTJe1z898OWj3Bg6kMNnUQi9RCtOMymODrH03l6ake0tT9E7AbPE
1O6Gk9ftQ3Qe+/KFaMHmOrm9bVVrL1LCujwWT5A+P3bL+h2dVb+4o3Uz7uwTCjV1
S9idzqxRhgVTWwpKp9mjKNqSxXdezP3ci60RKuQ1DjJOdWiZj1RouKle4UwPYxB6
TNPD3OHzlKVvdWL7Z2Ksaq4+ukcDdzz9nFegVsvPPpXmY1hf3t+v2tZzKbQUcE6z
7veu/wjQWtkbLehGTq/1B9A7ZCR3UbLIOBAOhGPwi0DnIwcVZa3LYVABk+l6y0N2
w2ERPWk4c3JKu4pEgZqrk+u0W36tk2BPCk26ZbetVXKLdH5nIKrXW94OAM07y+Wd
O1Wow6RHt9DyyT9yV7WTknV0iQH74wIDAQABoy4wLDATBgNVHSUEDDAKBggrBgEF
BQcDATAVBgNVHREEDjAMggpmb29iYXIuY29tMA0GCSqGSIb3DQEBCwUAA4ICAQDX
XTbVGa8DDcdNeP8SfFcAxsknEwR9uN9G1XpyTyPYq5236yLDM4RNBhY+4jFx+emo
DWIDALECZuECe4W+fyGTy9Xmgc4ycWcbjmnWlkrEEX0SFvj1wKJOT5C8MHzhi+eA
txxhCD3x//GM2KBGrXDG1nqd5hpNVgn1zyHbbO8LCG7Q0Vmdr8063gHzV2vWIJxH
W11X5zMgEco75ZGvtcfESma136CD2fWsIkL9qCtFV82s+eabdTaqbvD/GJVrRNI1
Ri/tgjQ0kMZ7tRIQjOi15ImxFYpF7yji3nXuZ/x1ZVDlYiUzFGYEo/14onZb93YS
+fu4HWk+szQAT7ve7eFMzjItJuaa5Ar64KJxNFySw2SmqvSkbudxw0J2+iIwI93a
bmTrau23pMss9cP/KJrxH0Salzp+gTw90IQcTQavcPPEwVJso5LPn6dxjuoliYTm
3y8gCSOFVvjoJcFS+JX3G5OVT0hAYBP9rL9zbJh89SI7q2po9jv+1oBj8ChYt6aU
CjBUFgc64EcYJ7qtV1kFRH636Bx856RwWwSaR067jSJ96HpRJJg2RnwzokGjd7gF
H5/FdIWlAiM3dm37qGdGu5pG+Mx4ro7vsa82vAFs2/8wxrlFY2e9fyUQxXnk2izG
O0s/LbPH1mDDN8PoQUq5g7mL6RVQy2xmwKxOj3oaDQ==
-----END CERTIFICATE-----`

const inner_key = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

func NewTls(crt string, key string, addr string) (ln net.Listener, err error) {
	//TODO check empty and give default cer.
	cer, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		utils.Logger.Info("未设置证书或证书格式化错误，使用内置证书替代。")
		// 没有指定cert key 使用默认
		cer, err = tls.X509KeyPair([]byte(inner_cert), []byte(inner_key))
		if err != nil {
			utils.Logger.Error(err.Error())
			return
		}
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err = tls.Listen("tcp", addr, config)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}

	return
}
