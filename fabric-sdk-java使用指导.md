### fabric-sdk-java使用指导
	本文档指导说明成功运行一个SDK测试用例，github上提供了使用说明，该指导文档以官方说明为基础做简单修改，不同版本使用方式可能有细微差别，仅供参考。
github地址：https://github.com/hyperledger/fabric-sdk-java

使用步骤如下
  
1.	下载fabric-SDK-java

2.	配置mvn

3.	pom.xml添加如下依赖：

```
<dependency>
        <groupId>org.hyperledger.fabric-sdk-java</groupId>
        <artifactId>fabric-sdk-java</artifactId>
        <version>1.0.1</version>
</dependency>
```
 
 
4.	修改配置文件

官方示例仅提供单元测试，所有相关配置均写死在代码中，正式开发时可以从network-config.yaml配置文件读取。此处仅说明如何修改以成功运行测试用例。

修改文件：src\test\java\org\hyperledger\fabric\sdk\testutils\TestConfig.java

1）	修改100行左右peer、orderer的IP、端口等信息即可运行。

2）	若TLS打开，需配置成grpcs，若TLS关闭需配置成grpc

3）	关于TLS配置在函数
private Properties getEndPointProperties(final String type, final String name) {}
默认使用的是server.crt，修改成ca.crt也可以运行，可根据需要自行配置。

	其他相关配置基本在该文件及后面的测试用例文件，可根据需要自行修改。
  
5.	证书拷贝

目录：src\test\fixture\sdkintegration\e2e-2Orgs\channel

删除原始证书，拷贝自己的证书。

6.	测试用例

官方提供了完成的测试用例 End2endIT.java

改用例包含了create channel、join channel、insta chaincode、instantitate chaincode、invoke、query等所有相关操作，该配置是按照fabric官方环境搭建指导文档配置。运行在我们的环境会有一些小问题，如果仅需要invoke和query功能，可参考下面自己实现的用例。
 
demo代码：[SendTx.java](https://github.com/berryjam/fabric-learning/blob/master/SendTx.java)

该测试用例测试的是example02的chaincode，若测试该chaincode修改最前面的一些参数信息即可运行。若测试其他chaincode，修改invoke()、query()函数的fcn及args即可。
完成修改后run setup()函数即可进行测试。

注：reconstructChannel()函数为初始化SDK的client及channel。
