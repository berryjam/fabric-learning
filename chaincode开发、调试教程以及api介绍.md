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





