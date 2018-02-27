# Chaincode教程

## 1. 概要

### 1.1 什么是Chaincode？

Chaincode是一个程序，用Go,node.js,java等编程语言编写，并实现了特定的接口（后面会详细介绍，分别为`Init`和`Invoke`）。Chaincode在一个安全的Docker容器中运行，该容器与背书peer进程隔离。Chaincode通过应用程序提交的事务来初始化和管理账本状态。

Chaincode通常处理区块链网络成员商定的业务逻辑，因此可以将其视为“智能合约”。由chaincode创建的状态仅限于该chaincode，不能由另一个chaincode直接访问。然而，在同一个区块链网络中，给定适当的权限，chaincode可以调用另一个chaincode来访问其状态。

### 1.2 两种角色



---

Chaincode

+ Init(stub ChaincodeStubInterface) pb.Response
+ Invoke(stub ChaincodeStubInterface) pb.Response

---

ChaincodeStubInterface

+ GetArgs() [][]byte
+ GetStringArgs() []string
+ GetFunctionAndParameters() (string, []string)
+ GetArgsSlice() ([]byte, error)
+ GetTxID() string
+ GetChannelID() string
+ InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response
+ GetState(key string) ([]byte, error)
+ PutState(key string, value []byte) error
+ DelState(key string) error
+ GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error)
+ GetStateByPartialCompositeKey(objectType string, keys []string) (StateQueryIteratorInterface, error)
+ CreateCompositeKey(objectType string, attributes []string) (string, error)
+ SplitCompositeKey(compositeKey string) (string. []string, error)
+ GetQueryResult(query string) (StateQueryIteratorInterface, error)
+ GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error)
+ GetCreator() ([]byte, error)
+ GetTransient() (map[string]byte, error)
+ GetBinding() ([]byte, error)
+ GetDecorations() map[string][]byte
+ GetSignedProposal() (*pb.SignedProposal, error)
+ GetTxTimestamp() (*timestamp.Timestamp, error)
+ SetEvent(name striing, payload []byte) error


---

CommonIteratorInterface

+ HasNext() bool
+ Close() error

---


StateQueryIteratorInterface

+ Next() (*queryresult.KV, error)


---


HistoryQueryIteratorInterface

+ Next() (*queryresult.KeyModification, error)


---


MockQueryIteratorInterface


# chaincode开发、调试入门教程与相关api分析

**Note.** 本篇主要内容分为两部分，分别为chaincode开发调试入门教程和chaincode api分析，希望能帮助更多开发者熟悉chaincode和编写出满足自身业务的chaincode应用。

## 1.chaincode开发、调试入门教程



## 2.chaincode api分析

### 2.1 Chaincode 



### 2.2 ChaincodeStubInterface 

这个接口能够被**可部署**的chaincode应用来读写它们区块账本，这里的可部署是指用户编写的链代码以及部分系统链代码，因为有些系统链代码是不能够部署的，这些不可部署的链代码就不能使用读写账本。





