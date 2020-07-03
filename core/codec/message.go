package codec

//默认消息格式版本
const DefaultMessageVersion	= "1.0"

/*
	消息类型模版定义
*/
const (
	//属性类型模版
	MessageTypeTemplate_Property 	= "thing.property.post"
	//事件类型模版
	MessageTypeTemplate_Event 		= "thing.event.%s.post"
	//服务调用类型模版
	MessageTypeTemplate_Service 	= "thing.service.%s.call"
)

/*
	元信息结构定义
*/

//服务调用元信息结构定义
type ServiceMetadata struct {
	//模型id
	ModelId 	string 		`json:"modelId"`
	//实体id
	EntityId 	string 		`json:"entityId"`
}

//模型属性或事件元信息结构定义
type ModelMetadata struct {
	//模型id
	ModelId 	string 		`json:"modelId"`
	//实体id
	EntityId 	string 		`json:"entityId"`
	//设备源
	Source		[]string 	`json:"source"`
	//采样时间戳
	EpochTime 	int64 		`json:"epochTime"`
}

/*
	MDMP消息Payload头结构定义
*/
type MdmpMsgHeader struct {
	//消息id
	ID			string 					`json:"id"`
	//协议版本号
	Version 	string 					`json:"version"`
	//消息类型
	Type 		string 					`json:"type"`
	//消息元信息
	Metadata 	interface{}				`json:"metadata"`
}

type MdmpMsgReplyHeader struct {
	//消息id
	ID			string 					`json:"id"`
	//协议版本号
	Version 	string 					`json:"version"`
	//状态码
	Code 		int 					`json:"code"`
}

/*
	服务调用消息结构定义
*/
type MdmpServiceCallMsg struct {
	MdmpMsgHeader
	//服务调用数据
	Params 		map[string]interface{} 	`json:"params"`
}

/*
	模型属性消息结构定义
*/
type ModelPropertyData struct {
	//属性值
	Value		interface{}				`json:"value"`
	//时间戳
	Time 		int64 					`json:"time"`
}

type MdmpPropertyMsg struct {
	MdmpMsgHeader
	//模型属性数据
	Params 		map[string]*ModelPropertyData 	`json:"params"`
}

/*
	模型事件消息结构定义
*/
type ModelEventData struct {
	//事件内容
	Value		map[string]interface{}	`json:"value"`
	//信息，可为空
	Message 	string 					`json:"message"`
	//级别，可为空
	Level 		string 					`json:"level"`
	//时间戳
	Time 		int64 					`json:"time"`
}

type MdmpEventMsg struct {
	MdmpMsgHeader
	//模型时间数据
	Params 		*ModelEventData			`json:"params"`
}