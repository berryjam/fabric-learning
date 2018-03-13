# Chaincode教程

-  [1. 概要](https://github.com/berryjam/fabric-learning/blob/master/chaincode%E6%95%99%E7%A8%8B.md#1-%E6%A6%82%E8%A6%81)

    - [1.1 什么是Chaincode?](https://github.com/berryjam/fabric-learning/blob/master/chaincode%E6%95%99%E7%A8%8B.md#11-%E4%BB%80%E4%B9%88%E6%98%AFchaincode)
  
    - [1.2 两种角色](https://github.com/berryjam/fabric-learning/blob/master/chaincode%E6%95%99%E7%A8%8B.md#12-%E4%B8%A4%E7%A7%8D%E8%A7%92%E8%89%B2)
  
- [2. chaincode开发者教程](https://github.com/berryjam/fabric-learning/blob/master/chaincode%E6%95%99%E7%A8%8B.md#2-chaincode%E5%BC%80%E5%8F%91%E8%80%85%E6%95%99%E7%A8%8B)

    - [2.1 chaincode API]()
    
    - [2.2 示例：“简单资产管理” chaincode]()
    
    - [2.3 安装Hyperledger Fabric示例]()
    
    - [2.4 下载fabric相关docker镜像]()
    
    - [2.5 终端1-启动示例网络]()
    
    - [2.6 终端2-编译&启动chaincode]()
    
    - [2.7 终端3-使用chaincode]()
    
    - [2.8 测试新chaincode]()
    
    - [2.9 chaincode加密]()
    
    - [2.10 管理go语言编写的chaincode外部依赖]()
    
- [3. chaincode运维者教程](https://github.com/berryjam/fabric-learning/blob/master/chaincode%E6%95%99%E7%A8%8B.md#3-chaincode%E8%BF%90%E7%BB%B4%E8%80%85%E6%95%99%E7%A8%8B)

- [4. 参考]()

## 1. 概要

### 1.1 什么是Chaincode？

Chaincode是一个程序，用Go,node.js,java等编程语言编写，并实现了特定的接口（后面会详细介绍，分别为`Init`和`Invoke`）。Chaincode在一个安全的Docker容器中运行，该容器与背书peer进程隔离。Chaincode通过应用程序提交的事务来初始化和管理账本状态。

Chaincode通常处理区块链网络成员商定的业务逻辑，因此可以将其视为“智能合约”。由chaincode创建的状态仅限于该chaincode，不能由另一个chaincode直接访问。然而，在同一个区块链网络中，给定适当的权限，chaincode可以调用另一个chaincode来访问其状态。

### 1.2 两种角色

我们可以从两种不同的角色来认识chaincode。一个是从应用程序开发人员的角度出发，应用开发者会开发一个名为[Chaincode for Developers](https://github.com/berryjam/fabric-learning/blob/master/chaincode%E6%95%99%E7%A8%8B.md#2-chaincode%E5%BC%80%E5%8F%91%E8%80%85%E6%95%99%E7%A8%8B)的区块链应用程序／解决方案；另一个是面向区块链网络运维人员[Chaincode for Operators](https://github.com/berryjam/fabric-learning/blob/master/chaincode%E6%95%99%E7%A8%8B.md#3-chaincode%E8%BF%90%E7%BB%B4%E8%80%85%E6%95%99%E7%A8%8B)，区块链网络运维人员负责管理区块链网络，并利用Hyperledger Fabric API来安装、实例化和升级chaincode，但很可能不会涉及chaincode应用程序的开发。

下面我们将分别从chaincode开发者和运维人员两方面对chaincode做一个较为详细的介绍，最后通过结合源码分析，加深对chaincode的理解。最后希望能帮助chaincode开发者能快速上手chaincode的开发，还有帮助chaincode运维人员能够保证chaincode能正常的运行。

## 2. chaincode开发者教程



## 3. chaincode运维者教程 

## 4. 参考

---


<div align="center">
<img src="https://github.com/berryjam/fabric-learning/blob/master/markdown_graph/chaincode-class-diagram.jpeg?raw=true">
</div>

<p align="center">
  <b>图 1 chaincode api类图</b><br>
</p>






