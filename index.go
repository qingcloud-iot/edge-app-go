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
	Type 			common.AppSdkRuntimeType
	//消息回调处理函数
	MessageCB   	common.AppSdkMessageCB
	//消息回调处理函数的用户自定义参数
	MessageParam  	interface{}
	//事件回调处理函数
	EventCB			common.AppSdkEventCB
	//事件回调处理函数的用户自定义参数
	EventParam    	interface{}
	//订阅的服务调用的id数组
	ServiceIds		[]string
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
	SendMessage(*common.AppSdkMessageData) error
}

func NewClient(opt *Options) (Client, error) {
	if opt == nil {
		return nil, errors.New("options is nil")
	}
	obj := core.NewAppCoreClient(opt.Type, opt.MessageCB, opt.MessageParam, opt.EventCB, opt.EventParam, opt.ServiceIds)
	return obj, nil
}
