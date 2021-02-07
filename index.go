package edge_app_go

import (
	"errors"
	"github.com/qingcloud-iot/edge-app-go/common"
	"github.com/qingcloud-iot/edge-app-go/core"
)

/*
	SDK初始化参数结构定义
*/
type Options struct {
	//应用类型
	Type 				common.AppSdkRuntimeType
	//模型消息回调处理函数
	MessageCB   		common.AppSdkMessageCB
	//模型消息回调处理函数的用户自定义参数
	MessageParam  		interface{}
	//SDK事件回调处理函数
	EventCB				common.AppSdkEventCB
	//SDK事件回调处理函数的用户自定义参数
	EventParam    		interface{}
	//订阅的边设备服务调用的id数组
	ServiceIds			[]string
	//订阅子设备消息模型ID数组
	EndpointThingIds 	[]string
}

/*
	SDK接口定义
*/
type Client interface {
	//初始化SDK
	Init() error
	//清除SDK
	Cleanup()
	//启动SDK
	Start() error
	//停止SDK
	Stop()
	//发送消息
	SendMessage(msgType common.AppSdkMessageType, payload []byte) error
	//获取边设备信息
	GetEdgeDeviceInfo() (*common.EdgeLocalInfo, error)
	//获取子设备信息列表
	GetEndpointInfos() ([]*common.EndpointInfo, error)
	//调用子设备服务调用
	CallEndpoint(thingId string, deviceId string, req *common.AppSdkMsgServiceCall) (*common.AppSdkMsgServiceReply, error)
}

func NewClient(opt *Options) (Client, error) {
	if opt == nil {
		return nil, errors.New("options is nil")
	}
	obj := core.NewAppCoreClient(opt.Type, opt.MessageCB, opt.MessageParam,
		opt.EventCB, opt.EventParam, opt.ServiceIds, opt.EndpointThingIds)
	return obj, nil
}
