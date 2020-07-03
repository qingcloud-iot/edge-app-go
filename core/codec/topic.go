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
)

/*
	topic模版定义
*/
const (
	/*
		订阅属性消息Topic模版 /edge/{appId}/thing/event/property/post
		{appId}：应用id
	*/
	topicTemplate_SubProperty 	= "/edge/%s/thing/property/base/post"

	/*
		发布属性消息Topic模版 /edge/{appId}/thing/event/property/control
		{appId}：应用id
	*/
	topicTemplate_PubProperty 	= "/edge/%s/thing/property/base/control"

	/*
		订阅事件消息Topic模版 /edge/{appId}/thing/event/{identifier}/post
		{appId}：应用id
		{identifier}：事件标识id
	*/
	topicTemplate_SubEvent 		= "/edge/%s/thing/event/%s/post"

	/*
		发布事件消息Topic模版 /edge/{appId}/thing/event/{identifier}/control
		{appId}：应用id
		{identifier}：事件标识id
	*/
	topicTemplate_PubEvent 		= "/edge/%s/thing/event/%s/control"

	/*
		发布服务调用Topic模版 /edge/{appId}/thing/service/{identifier}/call
		{appId}：应用id
		{identifier}：服务标识id
	*/
	topicTemplate_PubService 	= "/edge/%s/thing/service/%s/call"
)