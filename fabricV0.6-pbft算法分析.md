# fabric拜占庭容错算法分析

**Note：** fabric在v0.6分支实现了pbft算法，下面对其实现进行分析，以便能进一步掌握pbft算法以及了解如何在fabric实现共识算法插件，使得fabric支持不同的共识算法。


整个consensus模块的流程大致为：

## 1.客户端向某个peer节点发送执行链代码请求

1.1 客户端通过调用fabric的RESTful接口/chaincode**调用链代码**或者**部署链代码**，fabric在处理请求的时候（fabric/core/rest/rest_api.go.ProcessChaincode）再通过JSON RPC向peer节点发起执行事务请求，hyperledger/fabric/core/devops.go的Deplopy、invokeOrQuery方法，会调用peer.Impl（这个结构提供peer服务的实现）的ExecuteTransaction方法，如下面代码所示：

```
//ExecuteTransaction executes transactions decides to do execute in dev or prod mode
func (p *Impl) ExecuteTransaction(transaction *pb.Transaction) (response *pb.Response) {
        if p.isValidator {
                response = p.sendTransactionsToLocalEngine(transaction)
        } else {
                peerAddresses := p.discHelper.GetRandomNodes(1)
                response = p.SendTransactionsToPeer(peerAddresses[0], transaction)
        }
        return response
}
// hyperledger/fabric/core/peer/peer.go
```

1.2 peer节点在启动时，读取配置文件core.yaml的文件配置项peer.validator.enabled的值，peer根据这个值将自身设置为validator或者非validator。validator与非validator的区别在于：前者能够直接执行事务，而后者不直接执行事务而是通过gRPC的方式调用validator节点来执行事务（相当于转发事务），详细请参见SendTransactionsToPeer的实现，最终请求会定向到sendTransactionsToLocalEngine。重点分析sendTransactionsToLocalEngine方法。

1.3 sendTransactionsToLocalEngin方法会调用`p.engine.ProcessTransactionMsg`，`p.engine`为结构体EngineImpl，这是Engine接口实例，在启动peer时候创建。Engine这个接口用于管理peer网络的通讯和处理事务。EngineImpl的结构如下：

```
// EngineImpl implements a struct to hold consensus.Consenter, PeerEndpoint and MessageFan
type EngineImpl struct {
        consenter    consensus.Consenter // 每个共识插件都需要实现Consenter接口，包括RecvMsg方法和ExecutionConsumer接口的里函数（可以直接返回）
        helper       *Helper // 包含一些工具类方法，可以调用外部接口，比如获取网络信息，消息签名、验证，持久化一些对象等
        peerEndpoint *pb.PeerEndpoint
        consensusFan *util.MessageFan
}

// hyperledger/fabric/consensus/helper/engine.go
```

1.4 `ProcessTransactionMsg`的代码如下，可以看见链代码查询事务直接执行不需要进行共识，因为读取某个peer节点的账本不会影响自身以及其他peer节点账本，所以不需要共识来同步。而链代码调用和部署事务会影响到单个peer节点账本和状态，所以会调用共识插件的RecvMsg函数来保证各个peer节点的账本和状态一致。

```
// ProcessTransactionMsg processes a Message in context of a Transaction
func (eng *EngineImpl) ProcessTransactionMsg(msg *pb.Message, tx *pb.Transaction) (response *pb.Response) {
        //TODO: Do we always verify security, or can we supply a flag on the invoke ot this functions so to bypass check for locally generated transactions?
        if tx.Type == pb.Transaction_CHAINCODE_QUERY {
           // ... 
           result, _, err := chaincode.Execute(cxt, chaincode.GetChain(chaincode.DefaultChain), tx) // 直接执行查询事务，不需要共识
           // ...
        } else {
           // ...
           err := eng.consenter.RecvMsg(msg, eng.peerEndpoint.ID)  // 使用共识插件保证各个peer节点账本和状态保持一致
           if err != nil {
                    response = &pb.Response{Status: pb.Response_FAILURE, Msg: []byte(err.Error())}
           }
           // ...
// hyperledger/fabric/consensus/helper/engine.go
```

## 2.收到链代码执行请求的peer节点对链代码执行、部署事务进行共识

**Note：** fabric V0.6分支实现了两种公式算法NOOPS和PBFT，默认是NOOPS，peer节点在启动根据配置文件core.yaml文件配置项peer.validator.consensus.plugin选择采用哪种共识算法。

- NOOPS：用于开发和测试使用的插件，当一个validator节点收到一个事务消息时，会把消息转为共识消息，并会向所有节点广播共识消息。所有节点都会接收到这条共识消息，并执行消息里的事务。`这是一种比较朴素的共识方式，一旦因为网络或者其他原因，有些节点没收到广播消息，就会存在状态不一致问题。`

- PBFT：PBFT算法实现。简单地说当网络里的错误失效节点数量f与总的节点数量N满足关系N>3f时，PBFT算法也能保证各个节点的状态保持一致。当然还有更进一步的系统约束条件，后面会说明。



- obcBatch能够批量地对消息进行共识，提高pbft的共识效率，因为如果一条消息就进行一次共识，成本会很高。events.Manager整个事件管理器，最上层peer的操作会通过events.Manager.Queue()来输入事件，再由事件驱动pbftCore等结构体去完成整个共识过程。

- event.Timer是用于管理时间驱动的事件的接口，比golang timer多了一些特性：就算timer已经触发，但是只要event thread调用stop或者reset，那么timer触发的event就不会分发到event queue。event.timerImpl主要在一个loop方法里处理接收到的event,从startChan接收event，并把接收到的event发送到events.Manager.events通道。而events.Manager会在eventLoop()里循环处理接收到的事件。

- 新建obcBatch时：

    - 创建一个batchTimer定时器，根据配置设定batchTimeout等信息。obcBatch设置了batchTimer，每当出现timeOut后，会发送一个RequestBatch事件；
    
    - 创建一个pbftCore，并设置pbft.requestTimeout和pbft.nullRequestTimeout；
    
下图1是consensus包的类图，包含了主要结构、接口以及之间的关系。（图片有点小，可以点击放大查看：））

<div align="center">
<img src="https://github.com/berryjam/fabric-learning/blob/master/markdown_graph/consensus-class-diagram.png?raw=true">
</div>

<p align="center">
  <b>图 1 consensus包的类图</b><br>
</p>


