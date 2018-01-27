# fabric拜占庭容错算法分析

**Note：** fabric在v0.6分支实现了pbft算法，下面对其实现进行分析，以便能进一步掌握pbft算法以及了解如何在fabric实现共识算法插件，使得fabric支持不同的共识算法。


整个consensus模块的流程大致为：
- obcBatch是事件驱动，events.Manager整个事件管理器，最上层peer的操作会通过events.Manager.Queue()来输入事件，再由事件驱动pbftCore等结构体去完成整个共识过程。

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
