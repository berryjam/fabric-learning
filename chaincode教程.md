# Chaincode教程

-  [1.概要]()

    - [1.1 什么是Chaincode?]()
  
    - [1.2 两种角色]()
  

## 1. 概要

### 1.1 什么是Chaincode？

Chaincode是一个程序，用Go,node.js,java等编程语言编写，并实现了特定的接口（后面会详细介绍，分别为`Init`和`Invoke`）。Chaincode在一个安全的Docker容器中运行，该容器与背书peer进程隔离。Chaincode通过应用程序提交的事务来初始化和管理账本状态。

Chaincode通常处理区块链网络成员商定的业务逻辑，因此可以将其视为“智能合约”。由chaincode创建的状态仅限于该chaincode，不能由另一个chaincode直接访问。然而，在同一个区块链网络中，给定适当的权限，chaincode可以调用另一个chaincode来访问其状态。

### 1.2 两种角色

我们可以从两种不同的角色来认识chaincode。一个是从应用程序开发人员的角度出发，应用开发者会开发一个名为[Chaincode for Developers](https://github.com/berryjam/fabric-learning/blob/master/chaincode%E6%95%99%E7%A8%8B.md#2-chaincode%E5%BC%80%E5%8F%91%E8%80%85%E6%95%99%E7%A8%8B)的区块链应用程序／解决方案；另一个是面向区块链网络运维人员[Chaincode for Operators](https://github.com/berryjam/fabric-learning/blob/master/chaincode%E6%95%99%E7%A8%8B.md#3-chaincode%E8%BF%90%E7%BB%B4%E8%80%85%E6%95%99%E7%A8%8B)，区块链网络运维人员负责管理区块链网络，并利用Hyperledger Fabric API来安装、实例化和升级chaincode，但很可能不会涉及chaincode应用程序的开发。

下面我们将分别从chaincode开发者和运维人员两方面对chaincode做一个较为详细的介绍，最后通过结合源码分析，加深对chaincode的理解。最后希望能帮助chaincode开发者能快速上手chaincode的开发，还有帮助chaincode运维人员能够保证chaincode能正常的运行。

## 2. chaincode开发者教程



## 3. chaincode运维者教程 

---


<div align="center">
<img src="https://github.com/berryjam/fabric-learning/blob/master/markdown_graph/chaincode-class-diagram.jpeg?raw=true">
</div>

<p align="center">
  <b>图 1 chaincode api类图</b><br>
</p>

# chaincode开发、调试入门教程与相关api分析

**Note.** 本篇主要内容分为两部分，分别为chaincode开发调试入门教程和chaincode api分析，希望能帮助更多开发者熟悉chaincode和编写出满足自身业务的chaincode应用。

## 1.chaincode开发、调试入门教程



## 2.chaincode api分析

### 2.1 Chaincode 



### 2.2 ChaincodeStubInterface 

这个接口能够被**可部署**的chaincode应用来读写它们区块账本，这里的可部署是指用户编写的链代码以及部分系统链代码，因为有些系统链代码是不能够部署的，这些不可部署的链代码就不能使用读写账本。





