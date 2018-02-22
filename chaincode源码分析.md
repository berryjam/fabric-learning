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
