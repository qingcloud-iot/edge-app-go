package codec

/*
	topic类型定义
*/
const (
	//订阅属性类型
	TopicType_SubProperty 		= "TopicType_SubProperty"
	//发布属性类型
	TopicType_PubProperty 		= "TopicType_PubProperty"
	//订阅事件类型
	TopicType_SubEvent			= "TopicType_SubEvent"
	//发布事件类型
	TopicType_PubEvent			= "TopicType_PubEvent"
	//发布服务调用类型
	TopicType_PubService 		= "TopicType_PubService"
	//订阅服务调用类型
	TopicType_SubService		= "TopicType_SubService"
	//发布服务调用回应类型
	TopicType_PubServiceReply 	= "TopicType_PubServiceReply"
	//订阅服务调用回应类型
	TopicType_SubServiceReply 	= "TopicType_SubServiceReply"
)

/*
	topic模版定义V1：消息代理模式，通过AppControl服务代理应用消息
*/
const (
	/*
		订阅属性消息Topic模版 /edge/{appId}/thing/event/property/post
		{appId}：应用id
	*/
	topicTemplateV1_SubProperty 	= "/edge/%s/thing/property/base/post"

	/*
		发布属性消息Topic模版 /edge/{appId}/thing/event/property/control
		{appId}：应用id
	*/
	topicTemplateV1_PubProperty 	= "/edge/%s/thing/property/base/control"

	/*
		订阅事件消息Topic模版 /edge/{appId}/thing/event/{identifier}/post
		{appId}：应用id
		{identifier}：事件标识id
	*/
	topicTemplateV1_SubEvent 		= "/edge/%s/thing/event/%s/post"

	/*
		发布事件消息Topic模版 /edge/{appId}/thing/event/{identifier}/control
		{appId}：应用id
		{identifier}：事件标识id
	*/
	topicTemplateV1_PubEvent 		= "/edge/%s/thing/event/%s/control"

	/*
		发布服务调用Topic模版 /edge/{appId}/thing/service/{identifier}/call
		{appId}：应用id
		{identifier}：服务标识id
	*/
	topicTemplateV1_PubService 	= "/edge/%s/thing/service/%s/call"

	/*
		订阅服务调用Topic模版  /sys/{thingId}/{deviceId}/thing/service/{Identifier}/call
		{thingId}: 边设备模型id
		{deviceId}: 边设备id
		{Identifier}: 服务标识id
	*/
	topicTemplateV1_SubService 	= "/sys/%s/%s/thing/service/%s/call"

	/*
		发布服务调用回应Topic模版 /sys/{thingId}/{deviceId}/thing/service/{Identifier}/call_reply
	 	{thingId}: 边设备模型id
		{deviceId}: 边设备id
		{Identifier}: 服务标识id
	*/
	topicTemplateV1_PubServiceReply 	= "/sys/%s/%s/thing/service/%s/call_reply"

	/*
	订阅服务调用回应Topic模版 /sys/{thingId}/{deviceId}/thing/service/{Identifier}/call_reply
	{thingId}：模型id
	{deviceId}：设备id
	{Identifier}：服务标识id
*/
	topicTemplateV1_SubServiceReply = "/sys/%s/%s/thing/service/%s/call_reply"
)

/*
	topic模版定义V2：普通模式，直接使用平台消息格式通信
*/
const (
	/*
		订阅属性消息Topic模版 /sys/{thingId}/{deviceId}/thing/property/base/post
		{thingId}：模型id
		{deviceId}：设备id
	*/
	topicTemplateV2_SubProperty 	= "/sys/%s/%s/thing/property/base/post"

	/*
		发布属性消息Topic模版 /sys/{thingId}/{deviceId}/thing/property/base/post
		{thingId}：模型id
		{deviceId}：设备id
	*/
	topicTemplateV2_PubProperty 	= "/sys/%s/%s/thing/property/base/post"

	/*
		订阅事件消息Topic模版 /sys/{thingId}/{deviceId}/thing/event/{Identifier}/post
		{thingId}：模型id
		{deviceId}：设备id
		{Identifier}：事件标识id
	*/
	topicTemplateV2_SubEvent 		= "/sys/%s/%s/thing/event/%s/post"

	/*
		发布事件消息Topic模版  /sys/{thingId}/{deviceId}/thing/event/{Identifier}/post
		{thingId}：模型id
		{deviceId}：设备id
		{Identifier}：事件标识id
	*/
	topicTemplateV2_PubEvent 		= "/sys/%s/%s/thing/event/%s/post"

	/*
		发布服务调用Topic模版 /sys/{thingId}/{deviceId}/thing/service/{Identifier}/call
		{thingId}：模型id
		{deviceId}：设备id
		{Identifier}：服务标识id
	*/
	topicTemplateV2_PubService 		= "/sys/%s/%s/thing/service/%s/call"

	/*
		订阅服务调用Topic模版  /sys/{thingId}/{deviceId}/thing/service/{Identifier}/call
		{thingId}：模型id
		{deviceId}：设备id
		{Identifier}：服务标识id
	*/
	topicTemplateV2_SubService 		= "/sys/%s/%s/thing/service/%s/call"

	/*
		发布服务调用回应Topic模版 /sys/{thingId}/{deviceId}/thing/service/{Identifier}/call_reply
		{thingId}：模型id
		{deviceId}：设备id
		{Identifier}：服务标识id
	*/
	topicTemplateV2_PubServiceReply = "/sys/%s/%s/thing/service/%s/call_reply"

	/*
		订阅服务调用回应Topic模版 /sys/{thingId}/{deviceId}/thing/service/{Identifier}/call_reply
		{thingId}：模型id
		{deviceId}：设备id
		{Identifier}：服务标识id
	*/
	topicTemplateV2_SubServiceReply = "/sys/%s/%s/thing/service/%s/call_reply"
)