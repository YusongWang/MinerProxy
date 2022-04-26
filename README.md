### FIX 
为未来兼容性做处理
- TODO 矿池端要知道第一个包之后才可以链接发送。开发者抽水线程也是一样的。如果没有包链接上。不会建立长链接。同时如果没有机器在线要进行下线处理。
- TODO 动态修改配置文件中的抽水比例等
- TODO 确保抽水池关闭后重新连接不会影响当前矿工等。


#### done 
- 目前被频繁打开端口会频繁请求矿池。会被矿池拉黑。要等到第一个有效封包进入之后再打开矿池。

## 更新记录
### v0.0.1
#### 第四周. 目标 预计完成时间: 2022-04-30

##### TODO 
- 1. TODO 新增Web相关功能API(子进程守护模式。IPC交互。启动终止及重启功能)
- 2. TODO 适配ASIC机器
##### 完成 
- 


#### 第三周. 目标 预计完成时间: 2022-04-17

##### TODO 
- 4. TODO 新增Web相关功能API(子进程守护模式。IPC交互。启动终止及重启功能)
- 5. TODO 已延期。需要有RPC交互后修改 ----动态修改配置文件中的抽水比例等
- 6. TODO 适配ASIC矿机器

##### 完成
- 1. 多机器在线任务记录旷工唯一主键处理
- 2. deamon web watch dog. 读取配置文件。如果配置文件有变动。子进程通知父线程。watch dog 会重启子线程应用新的web端口.
- 3. deamon watch dog 监控所有server proxy 进程。掉线，重启。关闭等需求。
- 4. 上送矿工状态。
- 5. 修改链接线程为收到第一个包的时候再链接到矿池

#### 第二周. 目标 预计完成时间: 2022-04-10
##### TODO 
- TODO 多机器在线任务记录旷工唯一主键处理
- TODO 适配ASIC矿机器

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
