# 边缘应用**SDK**文档 - Go

### 概述

-------

这边文档包含以下几个方面内容：

- Go版本依赖
- **SDK**获取
- 接口功能
- 使用简介
- 示例介绍

### Go版本

-------

```sh
1.13以及以上
```

### **SDK**获取

-------

```sh
go get github.com/qingcloud-iot/edge-app-go
```

### 接口功能

--------

接口定义文件：[index.go](https://github.com/qingcloud-iot/edge-app-go/blob/master/index.go)

|      | 接口                                  | 功能描述                   |
| ---- | ------------------------------------- | -------------------------- |
|   1  | Init                                  | 初始化SDK                   |
|   2  | Cleanup                               | 清除SDK                   |
|   3  | Start                                 | 启动SDK                   |
|   4  | Stop                                  | 停止SDK                   |
|   5  | SendMessage                           | 发送消息                   |

### **SDK**使用简介

-------

- SDK初始化和启动

```sh
    //初始化配置参数
    options := &edge_app_go.Options{
        //应用类型
		Type: common.AppSdkRuntimeType_Docker,
		//接收消息回调函数
		MessageCB: onMessage,
		//回调用户自定义参数
		MessageParam: nil,
		//接收SDK事件通知回调，比如连接成功，连接断开等通知
		EventCB: onSdkEvent,
		//回调用户自定义参数
		EventParam: nil,
	}
	//创建SDK对象
	client, err := edge_app_go.NewClient(options)
	if err != nil {
        ...
	}
	//初始化SDK
	err = client.Init()
	if err != nil {
        ...
	}
	//启动SDK
	err = client.Start()
	if err != nil {
        ...
	} 
``` 

- 回调处理

```sh
func onMessage(msg *common.AppSdkMessageData, param interface{}) {
	//判断消息类型
	if msg.Type == common.AppSdkMessageType_Property {
	    //处理模型属性消息
	    ...
	    //反序列化消息
	    props := make([]*AppSdkMsgProperty, 0)
	    err := json.Unmarshal(msg.Payload, &props)
	    if err != nil {
		    ...        
	    }	    
	} else if msg.Type == common.AppSdkMessageType_Event {
	    //处理模型事件消息
	    ...
	    //反序列化消息
		evt := &AppSdkMsgEvent{}
    	err := json.Unmarshal(msg.Payload, evt)
    	if err != nil {
    		...
    	}    
	} else {
	    ...
	}
}
```

- 发送消息

1. 发送模型属性消息

```sh
propMsg := &common.AppSdkMessageData{
	Type: common.AppSdkMessageType_Property,
}
propData := make([]*common.AppSdkMsgProperty, 0)
tempProp := &common.AppSdkMsgProperty{
    //模型属性Identifier，需要跟平台的数据模型匹配
	Identifier: RANDOM_DATA_PROPERTY_ID,    
	Timestamp: time.Now().UnixNano() / 1e6,
	Value: strconv.Itoa(value),
}
propData = append(propData, tempProp)
propMsg.Payload, _ = json.Marshal(propData)
err := cli.SendMessage(propMsg)
if err != nil {
	...
}    
```

2. 发送模型事件消息

```sh
evtData := &common.AppSdkMsgEvent{
    //模型事件Identifier，需要跟平台的数据模型匹配
	Identifier: EVENT_ID,
	Timestamp: time.Now().UnixNano() / 1e6,
	Params: make(map[string]interface{}),
}
//模型事件参数的Identifier，需要跟平台的数据模型匹配
evtData.Params[EVENT_PARAM_DATA] = strconv.Itoa(value)
evtMsg := &common.AppSdkMessageData{
    Type: common.AppSdkMessageType_Event,
}
evtMsg.Payload, _ = json.Marshal(evtData)
err := cli.SendMessage(evtMsg)
if err != nil {
	...
}
```

3. 发送服务调用

```sh
srvData := &common.AppSdkMsgServiceCall{
    ////服务调用的Identifier，需要跟平台的数据模型匹配
	Identifier: SERVICE_ID,
	Params: make(map[string]interface{}),
}
//服务调用参数的Identifier，需要跟平台的数据模型匹配
srvData.Params[SERVICE_PARAM_ID] = resultValue
	srvMsg := &common.AppSdkMessageData{
	Type: common.AppSdkMessageType_ServiceCall,
}
srvMsg.Payload, _ = json.Marshal(srvData)
err := cli.SendMessage(srvMsg)
if err != nil {
	...
}
```


### **示例介绍** 

-------

[example/checker](https://github.com/qingcloud-iot/edge-app-go/tree/master/example/checker)





