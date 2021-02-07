package common

//消息处理回调定义
type AppSdkMessageCB func(*AppSdkMessageData, interface{})

//SDK事件处理回调定义
type AppSdkEventCB func(*AppSdkEventData, interface{})

/*
	应用类型枚举定义
*/
type AppSdkRuntimeType int32

const (
	//未知应用类型
	AppSdkRuntimeType_Unknown AppSdkRuntimeType = iota
	//Docker
	AppSdkRuntimeType_Docker
	//二进制
	AppSdkRuntimeType_Exec
)

type MessageModelType int32

/*
	消息数据类型和消息数据结构定义
*/
type AppSdkMessageType int32

const (
	//未知消息类型
	AppSdkMessageType_Unknown AppSdkMessageType = iota
	//属性消息
	AppSdkMessageType_Property
	//事件消息
	AppSdkMessageType_Event
	//服务调用消息
	AppSdkMessageType_ServiceCall
	//服务调用消息回应
	AppSdkMessageType_ServiceReply
)

//消息数据结构体
type AppSdkMessageData struct {
	/*
		消息类型，参考AppSdkMessageType枚举说明
	*/
	Type 			AppSdkMessageType				`json:"type"`
	/*
		物模型ID
	*/
	ThingId 		string 							`json:"thingId"`
	/*
		设备ID
	*/
	DeviceId 		string 							`json:"deviceId"`
	/*
		消息内容，为JSON字符串格式
		1. 当为AppSdkMessageType_Property类型时，Payload内容格式为[]*AppSdkMsgProperty通过JSON序列化之后的字符串，如下：
			`[{"identifier":"id_prop_01","timestamp":1593274999806,"value":"aaaaaa"},{"identifier":"id_prop_02","timestamp":1593274999806,"value":"bbbbbb"}]`
		2. 当为AppSdkMessageType_Event类型时，Payload内容格式为*AppSdkMsgEvent通过JSON序列花之后的字符串，如下:
			`{"identifier":"test_event_001","timestamp":1593274999806,"params":{"param1":"aaa","param2":20,"param3":"ccc"}}`
		3. 当为AppSdkMessageType_ServiceCall类型时，Payload内容格式为*AppSdkMsgServiceCall序列化之后的字符串，如下:
			`{"messageId":"40682013-308D-43DF-B2A3-819D5CDB08BD","identifier":"test_service_001","params":{"param1":"aaa","param2":20,"param3":"ccc"}}`
		4. 当为AppSdkMessageType_ServiceReply类型时，Payload内容格式为*AppSdkMsgServiceReply序列化之后的字符串，如下:
			`{"messageId":"40682013-308D-43DF-B2A3-819D5CDB08BD","identifier":"test_service_001","code":200,"params":{"param1":"aaa","param2":20,"param3":"ccc"}}`
		5. 其他类型，暂不支持，payload为nil
	*/
	Payload 		[]byte					`json:"payload"`
}

//属性消息结构体，AppSdkMessageType为AppSdkMessageType_Property时的payload
type AppSdkMsgProperty struct {
	Identifier 		string 					`json:"identifier"`
	Timestamp 		int64 					`json:"timestamp"`
	Value 			interface{}				`json:"value"`
}

//事件消息结构体，AppSdkMessageType为AppSdkMessageType_Event时的payload
type AppSdkMsgEvent struct {
	Identifier 		string 					`json:"identifier"`
	Timestamp 		int64 					`json:"timestamp"`
	Params 			map[string]interface{} 	`json:"params"`
}

//服务调用消息结构体，AppSdkMessageType为AppSdkMessageType_ServiceCall时的payload
type AppSdkMsgServiceCall struct {
	MessageId 		string 					`json:"messageId"`
	Identifier 		string 					`json:"identifier"`
	Params 			map[string]interface{} 	`json:"params"`
}

//服务调用消息结构体，AppSdkMessageType为AppSdkMessageType_ServiceReply时的payload
type AppSdkMsgServiceReply struct {
	MessageId 		string 					`json:"messageId"`
	Identifier 		string 					`json:"identifier"`
	Code			int32 					`json:"code"`
	Params 			map[string]interface{} 	`json:"params"`
}

/*
	SDK事件类型和事件数据结构定义
*/
type EventType int32

const (
	//未知事件类型
	EventType_Unknown EventType = iota
	//连接成功事件
	EventType_Connected
	//连接断开事件
	EventType_Disconnected
)

//SDK事件结构体
type AppSdkEventData struct {
	/*
		事件类型，参考EventType
	*/
	Type 		EventType
	/*
		事件数据
	*/
	Payload 	interface{}
}

