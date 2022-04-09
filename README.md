# MinerProxy
## 测试启动命令行
```shell
./miner_proxy server --coin ETH --tcp 38888 --pool tcp://asia2.ethermine.org:4444 --feepool tcp://asia2.ethermine.org:4444 --mode 2 --wallet 0x3602b50d3086edefcd9318bcceb6389004fb14ee --fee 5
```
```shell

#TEST
./miner_proxy server --coin ETH_TEST --tcp 38888 --pool ssl://api.wangyusong.com:8443 --feepool ssl://api.wangyusong.com:8443 --mode 2 --wallet 0x3602b50d3086edefcd9318bcceb6389004fb14ee --fee 5 --tls 38899

```


## 更新记录
### v0.0.1
#### 第二周. 目标 预计完成时间: 2022-04-10
##### TODO 
- 适配ASIC矿机器

Web界面相关
- 新增Web相关功能API(子进程守护模式。IPC交互。启动终止及重启功能)
- 已延期。需要有RPC交互后修改 ----动态修改配置文件中的抽水比例等

##### 完成
- 解析矿池难度。方便计算不同矿池抽水比例字段。


## FIX
- 清理过期任务防止内存爆炸

##### 完成
- 适配ETC

#### 第一周. 目标 预计 完成时间: 2022-04-02
#####  完成
- TCP SSL 适配
- 配置文件读取
- 抽水中转模式
- 普通中转模式
