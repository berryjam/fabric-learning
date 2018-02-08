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
    
      - `policies` 包含与所有peer节点可访问的链代码有关的策略，如背书策略。请注意，      
    
    
