# 2. 基本事务背书流程

下面我们从高层次的角度概述一个事务的请求流程。

**备注**：请注意，以下协议并不假定*所有事务都是确定性的，即允许非确定性事务*。

## 2.1. 客户端创建一个事务并将其发送给选择的背书peer节点

为了调用一个事务，客户端发送一个`PROPOSE`消息给它所选择的一组背书peer节点（可能不是同一时间发送 - 见2.1.2节和2.3节）。一个给定的`chaincodeID`的背书peer节点的集合经由peer节点被提供给客户端，该peer节点又从背书策略（参见第3节）知道背书peer节点的集合。例如，事务可以发送给给定`chaincodeID`的所有背书节点。也就是说，一些背书节点可能会离线，其他节点可能会反对，并选择不对事务背书。提交客户端尝试满足可用的背书节点的背书策略表达式。

在下文中，我们首先详细介绍`PROPOSE`消息格式，然后讨论提交客户端和背书节点之间可能的交互模式。

## 2.1.1. `PROPOSE`消息格式

`PROPOSE`消息的格式为`<PROPOSE,tx,[anchor]>`，其中`tx`是必要参数，`anchor`是可选参数。

- `tx=<clientID,chaincodeID,txPayload,timestamp,clientSig>`，其中：

    - `clientID`是提交客户端的ID；
    
    - `chaincodeID`是指事务所属的链代码的ID；

    - `txPayload`是包含提交的事务本身的有效数据载体；

    - `timestamp`是客户端维护的单调递增（对于每个新的事务）整数；
    
    - `clientSig`是tx其他字段上客户端的签名，即clientSig=Signature(clientID,chaincodeID,txPayload,timestamp)；
    
    `txPayload`的细节在调用事务和部署事务（如系统结构所述，部署事务其实是对系统链代码事务的调用）之间将有所不同。对于**调用事务**，`txPayload`由两个字段组成：
 
    - `txPayload = <operation, metadata>`，其中
    
      - `operation` 表示链代码操作（函数）和参数；
      
      - `metadata` 表示与调用相关的属性；
        

    对于**部署事务**来说，`txPayload`由三个字段组成：
    
    - `txPayload = <source, metadata, policies>`，其中
    
      - `source` 表示链代码的源码；
    
      - `metadata` 表示与链代码和应用程序相关的属性；
    
      - `policies` 包含与所有peer节点可访问的链代码有关的策略，如背书策略。请注意，背书策略在`deploy`事务不会与`txPayload`一起提供，但`deploy`的`txPayload`包含背书策略ID以及参数（请参阅第3节）。
      
- `anchor`包含*读版本依赖*，或者更具体地说，key-version对（即，`anchor`是`KxN`的一个子集），它将`PROPOSE`请求绑定或"anchors"到KVS指定版本的key。（请参阅1.2节）。如果客户端指定了`anchor`参数，则一个背书者仅在*读*取它本地KVS匹配`anchor`中相应key的版本号时，才会对一个事务进行背书（更多细节请参阅2.2节）。
    
`tx`的加密散列值被所有节点用作唯一事务标识符`tid`（即，`tid=HASH(tx)`）。客户端在内存中存储`tid`，并等待来自背书peer节点的响应。
    
## 2.1.1. 消息模式

客户端决定与背书peer节点的交互顺序。例如，客户端通常会发送`<PROPOSE, tx>`(即没有`anchor`参数)给一个单一背书peer节点，然后产生版本依赖（`anchor`)，客户端以后可以把所产生的`anchor`作为其`PROPOSE`消息的参数发送给其他背书节点。作为另外一个例子，客户端可以直接发送`<PROPOSE, tx>`(没有`anchor`)给所选的所有背书节点。不同模式的通讯都是可以的，并且客户端可以自由决定（见2.3节）。

## 2.2. 背书peer节点模拟事务并产生背书签名

在收到来自客户端的`<PROPOSE,tx,[anchor]>`时，背书peer节点`epID`首先验证客户端的签名`clientSig`，然后模拟执行一个事务。如果客户端指定了`anchor`，则只有在其本地KVS中的对应键的读取版本号（即，下面定义的`readset`）匹配`anchor`指定的版本号时，背书peer节点才会模拟执行事务。
    
    
模拟事务涉及通过背书peer节点调用事务引用的链代码(`chaincodeID`)以及复制背书peer节点在本地持有的状态副，以此来临时执行事务(`txPayload`)。

作为执行的结果，背书peer节点计算出*读取版本依赖性*(`readset`)和*状态更新*(`writeset`)，也称为DB语言中的*MVCC+postimage*。*MVCC*是多版本并发控制[Multiversion Concurrency control](https://en.wikipedia.org/wiki/Multiversion_concurrency_control)的缩写名词。*MVCC*一般用于数据库管理系统，为了解决存在多并发读写数据时可能出现数据不一致性问题。因为区块链系统中存在多个客户端请求执行事务，所以也会存在账本读写不一致问题。

回想一下，状态由key/value(k/v)对组成。所有k/v对都是版本化的，即每一个k/v对都包含了有序的版本化信息，每当某个key对应的value更新时，该版本信息就会增加。解释事务的peer节点会记录被链代码访问的所有k/v对，用于读取或写入，但peer节点尚未更新其状态。进一步来说：

- 给定一个在背书peer节点执行事务前，它的状态为`s`，对于事务读取的每个`k`，`(k,s(k).version)`pair被添加到到`readset`中。

- 另外，对由事务修改为新值`v'`的每个`k`，`(k,v')`pair被添加到`writeset`。或者，`v'`可能是新值与先前值(`s(k).value`)的增量。（使用增量可以减少`writeset`的大小）

如果客户端在`PROPOSE`消息中指定了`anchor`，则客户端指定的`anchor`必须等于模拟执行事务时由背书peer节点生成的`readset`。

然后，peer节点向内部转发`tran-proposal`（也可能是`tx`）到其事务背书的（peer）逻辑部分，称为**背书逻辑**。默认情况下，peer节点的背书逻辑会接受`tran-proposal`并简单地对`tran-proposal`进行签名。然而，背书逻辑可以解释任意的功能，例如，以`tran-proposal`和`tx`作为输入来与遗留系统交互，以达成是否认可交易的决定。

如果背书逻辑决定认可一个事务，它发送`<TRANSACTION-ENDORSED, tid, tran-proposal,epSig>`消息到提交客户端(`tx.clientID`)，其中：

- `tran-proposal := (epID,tid,chaincodeID,txContentBlob,readset,writeset)`，其中`txContentBlob`是(链代码/事务)特定的信息。目的是将`txContentBlob`用作`tx`的一些表示（例如`txContentBlob=tx.txPayload`）。

- `epSig`是背书peer节点对`tran-proposal`的签名。

否则，如果背书逻辑拒绝认可事务，则背书peer节点*可能*向客户端发送消息`(TRANSCATION-INVALID, tid, REJECTED)`。

请注意，背书peer节点在这一步不会改变其状态，在背书背景下


