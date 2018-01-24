# fabric拜占庭容错算法分析

**Note：** fabric在v0.6分支实现了pbft算法，下面对其实现进行分析，以便能进一步掌握pbft算法以及了解如何在fabric实现共识算法插件，使得fabric支持不同的共识算法。


整个consensus模块的流程大致为：
- obcBatch是事件驱动，events.Manager整个事件管理器，最上层peer的操作会通过events.Manager.Queue()来输入事件，再由事件驱动pbftCore等结构体去完成整个共识过程。

- event.Timer是用于管理时间驱动的事件的接口，比golang timer多了一些特性：就算timer已经触发，但是只要event thread调用stop或者reset，那么timer触发的event就不会分发到event queue。

- 新建obcBatch时，会创建一个batchTimer定时器，根据配置设定batchTimeout等信息。
