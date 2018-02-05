### 2.2 pbft实现

pbft算法的3段协议、VIEW-CHANGE协议、垃圾回收等等都可以看作是由一个个事件来驱动运行的。比如三段协议的pre-prepare阶段某primary节点收到客户端的请求这个事件后，就会广播pre-prepare消息；比如commit阶段，当一个节点收到来自于其他节点的commit消息事件后，可能会执行消息所捎带的事务。fabric在实现pbft的时候引入了`事件驱动模型`，代码在hyperledger/fabric/consensus/util/events/events.go。另外，为了提高共识效率，会批量处理pbft的消息，而不是一条一条处理。而批量处理是由Timer定时器触发，还有VIEW-CHANGE协议也会用到`Timer定时器`。当backup节点等待执行请求超时会发送一个VIEW-CHANGE消息，fabric实现了一个Timer定时器。之所以单独介绍Event模型和Timer定时器，是因为要想完全看懂pbft的实现，就必须理解它的事件流以及Timer定时器。

#### 2.2.1 Event模型

```

```

#### 2.2.2 Timer定时器

