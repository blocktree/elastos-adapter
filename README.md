# elastos-adapter

elastos-adapter适配了openwallet.AssetsAdapter接口，给应用提供了底层的区块链协议支持。

## 如何测试

openwtester包下的测试用例已经集成了openwallet钱包体系，创建conf文件，新建ELA.ini文件，编辑如下内容：

```ini


# node api url, if RPC Server Type = 1, use bitbay insight-api
serverAPI = "http://ip:port"
# use fixed fee or not
useFixedFee = false
# fixed fee
fixedFee = "0.005"
# Cache data file directory, default = "", current directory: ./data
dataDir = ""
```
