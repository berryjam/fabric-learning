### fabric-sdk-java使用指导
	本文档指导说明成功运行一个SDK测试用例，支持app通过sdk直连公有云的bcs实例，github上提供了使用说明，该指导文档以官方说明为基础做简单修改，不同版本使用方式可能有细微差别，仅供参考。
github地址：https://github.com/hyperledger/fabric-sdk-java

使用步骤如下

0.      在公有云上面已经订购了bcs，并且保证bcs实例所属集群的公网ip能够ping通，ip点击"通道管理"-"查看节点"来获取，如下图。如果不能ping通，请在安全组里面[添加ICMP规则](https://support.huaweicloud.com/usermanual-vpn/zh-cn_topic_0035557721.html)。

<div align="center">
<img src="https://github.com/berryjam/fabric-learning/blob/master/sdk_usage_pic/cluster_ip.png?raw=true">
</div>

1.	下载fabric-SDK-java，```git clone git@github.com:hyperledger/fabric-sdk-java.git```

2.	配置mvn
 
3.	修改配置文件

官方示例仅提供单元测试，所有相关配置均写死在代码中，正式开发时可以从network-config.yaml配置文件读取。此处仅说明如何修改以成功运行测试用例。

修改文件：src\test\java\org\hyperledger\fabric\sdk\testutils\TestConfig.java,具体修改内容可参考[TestConfig.java](https://github.com/berryjam/fabric-learning/blob/master/TestConfig.java)

1）	修改100行左右peer、orderer的IP、port、mspid等信息即可运行。

2）	当前公有云的环境只支持grcps方式，ip和port可以通过下载SDK配置文件来确定，ip为bcs实例所属的集群公网ip，而peer的port一般从30605开始。具体可参考下图以及相关字段说明。

<div align="center">
<img src="https://github.com/berryjam/fabric-learning/blob/master/sdk_usage_pic/testconfig_update.png?raw=true">
</div>

- peerOrg1.mspid：ea24fef7f9427f8086859fad278c7748e316b24cMSP（根据sdk配置内容修改，对应sdk的yaml文件的mspid）；

- peerOrg1.peer_locations："peer0.org1.example.com@grpcs://" + LOCALHOST + ":30605, peer1.org1.example.com@grpcs://" + LOCALHOST + ":30606" (**注意：LOCALHOST需要改为bcs所属集群公网ip，另外必须使用grpcs，具体端口号参考sdk配置yaml文件的url**) ```如：private static final String LOCALHOST = "49.4.14.76" ```；

- peerOrg1.orderer_locations："orderer.example.com@grpcs://" + LOCALHOST + ":30805"，这里orderer端口号默认是30805，具体以yaml文件的orderers的url里的端口为准，url一般为 grpcs://orderer-xxx-0.orderer-xxx.default.svc.cluster.local:30805；

- peerOrg1.eventhub_locations："peer0.org1.example.com@grpcs://" + LOCALHOST + ":30705,peer1.org1.example.com@grpcs://" + LOCALHOST + ":30706"，这里的evenHub端口一般也是从30705开始，具体端口号以yaml文件的eventUrl为准；

3）	关于TLS配置在函数
private Properties getEndPointProperties(final String type, final String name) {}
把server.crt，修改成ca.crt。

	其他相关配置基本在该文件及后面的测试用例文件，可根据需要自行修改。
  
4.	证书拷贝

目录：src\test\fixture\sdkintegration\e2e-2Orgs\channel

删除原始证书，拷贝自己的证书。

5.      把peer和orderer的根证书安装到jdk，否则会报以下错误。其中peer和orderer的根证书在tls/ca.crt,如果使用intellij运行，那么请参考[intellij安装证书](https://intellij-support.jetbrains.com/hc/en-us/community/posts/115000094584-IDEA-Ultimate-2016-3-4-throwing-unable-to-find-valid-certification-path-to-requested-target-when-trying-to-refresh-gradle)把tls/ca.crt安装到对应的jdk目录下面。如果使用命令行方式运行，请参考[jdk安装证书](https://blog.csdn.net/wn_hello/article/details/71600988)。

```
unable to find valid certification path to requested target
```


6.	测试用例

官方提供了完成的测试用例 End2endIT.java

改用例包含了create channel、join channel、insta chaincode、instantitate chaincode、invoke、query等所有相关操作，该配置是按照fabric官方环境搭建指导文档配置。运行在我们的环境会有一些小问题，如果仅需要invoke和query功能，可参考下面自己实现的用例。
 
demo代码：[SendTx.java](https://github.com/berryjam/fabric-learning/blob/master/SendTx.java)，把SendTx.java文件放到src\test\java\org\hyperledger\fabric\sdkintegration下面。

该测试用例测试的是example02的chaincode，若测试该chaincode修改最前面的一些参数信息即可运行。若测试其他chaincode，修改invoke()、query()函数的fcn及args即可。
完成修改后run setup()函数即可进行测试。

注：reconstructChannel()函数为初始化SDK的client及channel。

### 常见问题

1.```unable to find valid certification path to requested target```，jdk缺少bcs实例节点的根证书，需要把tls/ca.crt安装到jdk里。

2.```java.lang.RuntimeException: Missing cert file for: peer-ea24fef7f9427f8086859fad278c7748e316b24c-0.peer-ea24fef7f9427f8086859fad278c7748e316b24c.default.svc.cluster.local. Could not find at location: fabric-sdk-java/src/test/fixture/sdkintegration/e2e-2Orgs/v1.0/crypto-config/peerOrganizations/peer-ea24fef7f9427f8086859fad278c7748e316b24c.default.svc.cluster.local/peers/peer-ea24fef7f9427f8086859fad278c7748e316b24c-0.peer-ea24fef7f9427f8086859fad278c7748e316b24c.default.svc.cluster.local/tls/ca.crt``` ，sdk缺少证书，需要在创建相应的目录，并把msp、tls的证书拷贝到所创建的目录。

3.```First received frame was not SETTINGS. Hex dump for first 5 bytes ```，一般是因为TestConfig.java的配置不对，没有把grpc设置为grpcs，或者ip、port与sdk的yaml文件的ip、port不一致，解决方式参考上面第3节。

