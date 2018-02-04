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

- NOOPS：用于开发和测试使用的插件，当一个validator节点收到一个事务消息时，会把消息转为共识消息，并会向所有节点广播共识消息。一般情况下，所有节点都会接收到这条共识消息，并执行消息里的事务。`这是一种比较朴素的共识方式，一旦因为网络或者其他原因，有些节点没收到广播消息，就会存在状态不一致问题，所以不只用于开发和测试。`

- PBFT：PBFT算法实现。简单地说当网络里的错误失效节点数量f与总的节点数量N满足关系N>3f时，PBFT算法也能保证各个节点的状态保持一致。但是实现PBFT算法的需要满足以下的约束条件，**所以在选择共识算法时要对系统进行全面评估，基于系统自身情况选择，不能盲目选择。**：

    - 系统假设为异步分布式，通过网络传输的消息可能丢失、延迟、重复或者乱序。**作者假设节点的失效必须是独立发生的，也就是说代码、操作系统和管理员密码这些东西在各个节点上是不一样的**。
   
    - 使用了加密技术来防止欺骗攻击和重播攻击，以及检测被破坏的消息。**并且假设所有的节点都知道其他节点的公钥以进行签名验证。**

    - 可能存在多个失效、通讯存在延迟的节点,但是延迟节点不会无限期的被延迟，而且恶意攻击者算力有限无法破解加密算法。

### 2.1 pbft算法简介

**Note.** 下面pbft算法的介绍参考自[梧桐树博客](http://wutongtree.github.io/hyperledger/consensus)

#### 2.1.1 3阶段协议

在分析fabric-pbft的源码前，先对pbft算法流程做一个简单的描述。图1是pbft算法的三段协议过程：

<div align="center">
<img src="https://github.com/berryjam/fabric-learning/blob/master/markdown_graph/3-phase-protocol.jpg?raw=true">
</div>

<p align="center">
  <b>图 1 pbft算法三段协议过程</b><br>
</p>

从primary收到消息开始，每个消息都会有view的编号，每个节点都会检查是否和自己的view是相同的，代表是哪个节点发送出来的消息，源头在哪里，client收到消息也会检查该请求返回的所有消息是否是相同的view。如果过程中发现view不相同，消息就不会被处理。除了检查view之外，每个节点收到消息的时候都会检查对应的序列号n是否匹配，还会检查相同view和n的PRE-PREPARE、PREPARE消息是否匹配，从协议的连续性上提供了一定程度的安全。

每个节点收到其他节点发送的消息，能够验证其签名确认发送来源，但并不能确认发送节点是否伪造了消息，PBFT采用的办法就是数数，看有多少节点发送了相同的消息，在有问题的节点数有限的情况下，就能判断哪些节点发送的消息是真实的。REQUEST和PRE-PREPARE阶段还不涉及到消息的真实性，只是独立的生成或者确认view和序列号n，所以收到消息判断来源后就广播出去了。PREPARE阶段开始会汇总消息，通过数数判断消息的真实性。PREPARE消息是收到PRE-PREPARE消息的节点发送出来的，primary收到REQUEST消息后不会给自己发送PRE-PREPARE消息，也不会发送PRE-PREPARE消息，所以一个节点收到的消息数满足2f+1-1=2f个就能满足没问题的节点数比有问题节点多了（包括自身节点）。COMMIT阶段primary节点也会在收到PREPARE消息后发送COMMIT消息，所以收到的消息数满足2f+1个就能满足没问题的节点数比有问题节点多了（包括自身节点）。

PRE-PREPARE和PREPARE阶段保证了所有正常的节点对请求的处理顺序达成一致，它能够保证如果 PREPARE(m, v, n, i) 是真的话，PREPARE(m’, v, n, j) 就一定是假的，其中j是任意一个正常节点的编号，只要 D(m) != D(m’)。因为如果有3f+1个节点，至少有f+1个正常的节点发送了PRE-PREPARE和PREPARE消息，所以如果PREPARE(m’, v, n, j) 是真的话，这些节点中就至少有一个节点发了不同的PRE-PREPARE或者PREPARE消息，这和它是正常的节点不一致。当然，还有一个假设是安全强度是足够的，能够保证m != m’时，D(m) != D(m’)，D(m) 是消息m的摘要。

确定好了每个请求的处理顺序，怎么能保证按照顺序执行呢？网络消息都是无序到达的，每个节点达成一致的顺序也是不一样的，有可能在某个节点上n比n-1先达成一致。其实每个节点都会把PRE-PREPARE、PREPARE和COMMIT消息缓存起来，它们都会有一个状态来标识现在处理的情况，然后再按顺序处理。而且序列号n在不同view中也是连续的，所以n-1处理完了，处理n就好了。

#### 2.1.2 VIEW-CHANGE协议

<div align="center">
<img src="https://github.com/berryjam/fabric-learning/blob/master/markdown_graph/view-change.jpg?raw=true">
</div>

<p align="center">
  <b>图 2 VIEW-CHANGE协议过程</b><br>
</p>

上图是发生VIEW-CHANGE的一种情况，就是节点正常收到PRE-PREPARE消息以后都会启动一个定时器，如果在设置的时间内都没有收到回复，就会触发VIEW-CHANGE，该节点就不会再接收除CHECKPOINT 、VIEW-CHANGE和NEW-VIEW等消息外的其他消息了。NEW-VIEW是由新一轮的primary节点发送的，O是不包含捎带的REQUEST的PRE-PREPARE消息集合，计算方法如下：

- primary节点确定V中最新的稳定检查点序列号min-s和PRE-PREPARE消息中最大的序列号max-s；

- 对min-s和max-s之间每个序列号n都生成一个PRE-PREPARE消息。这可能有两种情况：

    - P的VIEW-CHANGE消息中至少存在一个集合，序列号是n；
    
    - 不存在上面的集合。
    
    第一种情况，会生成新的PRE-PREPARE消息<PRE-PREPARE, v+1, n, d>𝞂p，其中n是V中最大的v序列号，d是对应的PRE-PREPARE消息摘要。第二情况，PRE-PREPARE消息的d是特殊的空消息摘要。

primary节点发送完NEW-VIEW消息并记录到日志中就切换到v+1的view中，开始接收所有的消息了。其他节点也在收到NEW-VIEW消息后需要验证签名是否正确，还要验证O消息的正确性，都没问题就记录到日志中，广播完O中的PRE-PREPARE消息后就切换到v+1的view中，VIEW-CHANGE就算完成了。

#### 2.1.3 垃圾回收

每个节点都会把每条消息保存下来，除非它确认这个请求至少被f+1个正常节点处理过，而且还要能在VIEW-CHANGE中证明这一点。另外，如果一些节点错过了其他的正常节点都丢掉的消息，它需要传输部分或者全部的服务状态来保存同步。所以节点需要证明自己的状态是正确的。

如果每个操作完成都收集证据证明自己的状态没有问题成本就太高了。实际的做法可以是周期性的，比如请求的序号是100的倍数时。这种请求执行完的状态就叫一个检查点，验证过的检查点叫稳定检查点。每个节点维护了多个状态，最新的稳定检查点、多个不稳定的检查点和当前状态。

验证一个检查点的过程如下：

- 节点i生成一个检查点，广播<CHECKPOINT, n, d, i>𝞂i给其他的节点；

- 每个节点都检查自己的日志，如果有2f+1个序列号为n，消息摘要d相同的不同节点发送过来的CHECKPOINT消息，就是稳定检查点的证据；

确认了最新的稳定检查点，就可以把之前的检查点和检查点消息都删掉了，还可以删掉序列号小于n的所有PRE-PREPARE、PREPARE、COMMIT消息，减少占用的空间。

#### 2.1.4 一些优化措施

PBFT协议里提了几种优化措施：

- 减少通信：

    - 尽量避免发送大量的回复消息，client可以指定一个节点来发送回复消息，其他节点就只需要回复消息的摘要就可以了，这能在减少带宽和CPU开销的情况下验证结果的正确性；
    
    - 调用操作步骤从5步减少到4步。正常的调用需要经过REQUEST、PRE-PREPARE、PREPARE、COMMIT、REPLY等5步，节点可以在PREPARE后就处理消息，然后把执行结果发送给client，如果有2f+1个相同结果的消息，请求就结束了，否则还是正常的5步，出现异常的话就回退状态；
    
    - 提升只读操作的效率。节点只要能确认操作是正确而且是只读的，就可以立即执行，等待状态提交以后就回复给client；
    
- 节点采用签名来验证消息，实际使用的时候可以这么用：

    - 公钥签名：主要是VIEW_CHANGE、NEW_VIEW消息的时候用；
    
    - MAC：其他地方的消息传输都是这种方法，这样能减少性能瓶颈。MAC消息本来是不能验证消息的真实性，但是论文作者提供了一个办法来绕过这个问题，这会用一些规则，比如两个正常节点相同的v和n，请求也是一样的。

#### 2.1.5 一些优化措施

协议里面只介绍了主要的流程，很多实现的部分并没有说明，比如每个节点收到VIEW-CHANGE后怎么处理，MAC协议的共享密钥怎么分配，如果应对DDos攻击等等。

### 2.2 pbft实现

回到1.4节，当peer节点执行的是链代码调用或者部署事务时，需要进行共识，`err := eng.consenter.RecvMsg(msg, eng.peerEndpoint.ID)`。

- obcBatch能够批量地对消息进行共识，提高pbft的共识效率，因为如果一条消息就进行一次共识，成本会很高。events.Manager整个事件管理器，最上层peer的操作会通过events.Manager.Queue()来输入事件，再由事件驱动pbftCore等结构体去完成整个共识过程。

- event.Timer是用于管理时间驱动的事件的接口，比golang timer多了一些特性：就算timer已经触发，但是只要event thread调用stop或者reset，那么timer触发的event就不会分发到event queue。event.timerImpl主要在一个loop方法里处理接收到的event,从startChan接收event，并把接收到的event发送到events.Manager.events通道。而events.Manager会在eventLoop()里循环处理接收到的事件。

- 新建obcBatch时：

    - 创建一个batchTimer定时器，根据配置设定batchTimeout等信息。obcBatch设置了batchTimer，每当出现timeOut后，会发送一个RequestBatch事件；
    
    - 创建一个pbftCore，并设置pbft.requestTimeout和pbft.nullRequestTimeout；
    
下图2是consensus包的类图，包含了主要结构、接口以及之间的关系。（**图片有点小，可以点击放大查看：）**）

<div align="center">
<img src="https://github.com/berryjam/fabric-learning/blob/master/markdown_graph/consensus-class-diagram.png?raw=true">
</div>

<p align="center">
  <b>图 2 consensus包的类图</b><br>
</p>
