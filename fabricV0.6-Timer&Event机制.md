### 2.2 pbft实现

pbft算法的3段协议、VIEW-CHANGE协议、垃圾回收等等都可以看作是由一个个事件来驱动运行的。比如三段协议的pre-prepare阶段某primary节点收到客户端的请求这个事件后，就会广播pre-prepare消息；比如commit阶段，当一个节点收到来自于其他节点的commit消息事件后，可能会执行消息所捎带的事务。fabric在实现pbft的时候引入了`事件驱动模型`，代码在hyperledger/fabric/consensus/util/events/events.go。另外，为了提高共识效率，会批量处理pbft的消息，而不是一条一条处理。而批量处理是由Timer定时器触发，还有VIEW-CHANGE协议也会用到`Timer定时器`。当backup节点等待执行请求超时会发送一个VIEW-CHANGE消息，fabric实现了一个Timer定时器。之所以单独介绍Event模型和Timer定时器，是因为要想完全看懂pbft的实现，就必须理解它的事件流以及Timer定时器。

#### 2.2.1 Event模型

下面是事件管理器，Event的主要接口：

```
type Manager interface {
        Inject(Event)         // A temporary interface to allow the event manager thread to skip the queue
        Queue() chan<- Event  // Get a write-only reference to the queue, to submit events
        SetReceiver(Receiver) // Set the target to route events to
        Start()               // Starts the Manager thread TODO, these thread management things should probably go away
        Halt()                // Stops the Manager thread
}
```

事件管理器用于来管理事件，一般需要管理多个事件并且按事件接收的先后顺序来处理。因此需要有一个队列来存储事件，Queue()接口返回一个类型为Event的channel，用于存储事件。之所以使用channel，是因为Start()方法会启动一个goroutine循环处理接收到的事件，通过channel能够保证只有接收到事件才会处理，不用每时每刻循环检查队列去执行事件，浪费CPU性能。除了接收事件，还要能够处理事件。因此SetRecevier(Recevier)需要设置事件管理器的实际处理者，Recevier接口需要实现ProcessEvent(Event) Event方法。而obcBatch实现了这个方法，比如在处理一个committedEvent后会返回一个execDoneEvent，prepare消息又通过Queue()放到channel，在下一次的事件处理就会执行execDoneEvent，都是事件驱动的，符合pbft的算法模型。Start()方法会启动一个循环处理事件的goroutine：

```
// Start creates the go routine necessary to deliver events
func (em *managerImpl) Start() {
        go em.eventLoop()
}

// eventLoop is where the event thread loops, delivering events
func (em *managerImpl) eventLoop() {
        for {
                select {
                case next := <-em.events:
                        em.Inject(next)
                case <-em.exit:
                        logger.Debug("eventLoop told to exit")
                        return
                }
        }
}
```

eventLoop()方法会不断从事件队列channel取出事件，再通过Inject（Event）方法调用receiver来处理取出的事件。

```
// SendEvent performs the event loop on a receiver to completion
func SendEvent(receiver Receiver, event Event) {
        next := event
        for {
                // If an event returns something non-nil, then process it as a new event
                next = receiver.ProcessEvent(next)
                if next == nil {
                        break
                }
        }
}

// Inject can only safely be called by the managerImpl thread itself, it skips the queue
func (em *managerImpl) Inject(event Event) {
        if em.receiver != nil {
                SendEvent(em.receiver, event)
        }
}
```

Halt()方法用于停止循环处理事件。


#### 2.2.2 Timer定时器

之前提到过pbft里面会用到Timer定时器，比如backup只有在等待执行request超时的时候才会广播VIEW-CHANGE消息。下面是Timer接口：

```
type Timer interface {
        SoftReset(duration time.Duration, event Event) // start a new countdown, only if one is not already started
        Reset(duration time.Duration, event Event)     // start a new countdown, clear any pending events
        Stop()                                         // stop the countdown, clear any pending events
        Halt()                                         // Stops the Timer thread
}
```

SoftReset(time.Duration,Event)和Reset(time.Duration,Event)方法都会重新启动一个定时器，当启动时间超过duration就会处理event事件。这两个定时方法的区别是前者会先判断是否已经启动过定时器，如果是的话就忽略，否则才会启动；而后者会强制重置定时器。在Event模型已经描述过事件管理器处理event事件的流程，而Timer对象在实例化的过程中会设置Manager，从而达到定时处理Event的目的。

```
// newTimer creates a new instance of timerImpl
func newTimerImpl(manager Manager) Timer {
        et := &timerImpl{
                startChan: make(chan *timerStart),
                stopChan:  make(chan struct{}),
                threaded:  threaded{make(chan struct{})},
                manager:   manager, // 设置事件管理器
        }
        go et.loop() // 循环处理事件
        return et
}
```

#### 2.2.3 pbft共识代码

fabric V0.6分支的pbft公式算法代码都在位于文件夹consensus，consensus文件夹包含了controller、executor、helper、noops、pbft、util几个模块。

其中consensus.go包含了算法插件需要实现的RecvMsg()接口以及fabric外部提供给算法调用的接口，如执行管理账本状态的InvalidateState()、ValidateState()接口。

回顾1.4节，当peer节点执行调用链代码或者部署链代码的事务的时候，需要使用共识插件RecvMsg接口`err := eng.consenter.RecvMsg(msg, eng.peerEndpoint.ID)`对各个peer节点进行共识。接下来看pbft的RecvMsg的实现，如下：

```
// RecvMsg is called by the stack when a new message is received
func (eer *externalEventReceiver) RecvMsg(ocMsg *pb.Message, senderHandle *pb.PeerID) error {
        eer.manager.Queue() <- batchMessageEvent{
                msg:    ocMsg,
                sender: senderHandle,
        }
        return nil
}
```

如第2.2.1节Event模型所述，共识插件就会在循环等待接收Event事件，调用RecvMsg会向事件管理器EventManager传入一个batchMesageEvent，这个事件会捎带了从peer节点传进来的事务消息ocMsg，再通过receiver来处理接收到的Event事件。而pbft算法插件的recevier是obcBatch，能够批量处理共识消息。下面接着分析obcBatch是如何处理batchMessageEvent的：

```
// allow the primary to send a batch when the timer expires
func (op *obcBatch) ProcessEvent(event events.Event) events.Event {
        logger.Debugf("Replica %d batch main thread looping", op.pbft.id)
        switch et := event.(type) {  // 根据消息的反射类型来判断消息类型
        case batchMessageEvent:
                ocMsg := et
                return op.processMessage(ocMsg.msg, ocMsg.sender)  // ocMsg的消息类型仍为链代码事务类型
        case executedEvent:
                op.stack.Commit(nil, et.tag.([]byte))
        case committedEvent:
                logger.Debugf("Replica %d received committedEvent", op.pbft.id)
                return execDoneEvent{}
        // ...       
}
```
当接收到的是batchMessageEvent会调用processMessage来处理，并返回另外一种Event。接下来分析processMessage：

```
func (op *obcBatch) processMessage(ocMsg *pb.Message, senderHandle *pb.PeerID) events.Event {
        if ocMsg.Type == pb.Message_CHAIN_TRANSACTION {
                req := op.txToReq(ocMsg.Payload) // 这是pbft
                return op.submitToLeader(req)
        }
        
        // ....
}
```

##### **2.2.3.1 插件接口**

