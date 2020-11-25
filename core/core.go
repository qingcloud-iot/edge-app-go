package core

import (
	"errors"
	"fmt"
	"github.com/qingcloud-iot/edge-app-go/common"
	"github.com/qingcloud-iot/edge-app-go/core/codec"
	"github.com/qingcloud-iot/edge-app-go/core/config"
	"github.com/qingcloud-iot/edge-app-go/core/mqtt"
)

func NewAppCoreClient(appType common.AppSdkRuntimeType, msgCB common.AppSdkMessageCB, msgParam interface{},
						evtCB common.AppSdkEventCB, evtParam interface{}, srvIds []string) *AppCoreClient {
	return &AppCoreClient{
		appType: 		appType,
		messageCB: 		msgCB,
		messageParam: 	msgParam,
		eventCB: 		evtCB,
		eventParam:  	evtParam,
		serviceIds: 	srvIds,
	}
}

type AppCoreClient struct {
	//Runtime类型
	appType 		common.AppSdkRuntimeType
	//消息回调处理函数
	messageCB   	common.AppSdkMessageCB
	//消息回调处理函数的用户自定义参数
	messageParam  	interface{}
	//事件回调处理函数
	eventCB			common.AppSdkEventCB
	//事件回调处理函数的用户自定义参数
	eventParam    	interface{}
	//服务调用的id数组
	serviceIds 		[]string
	//mqtt协议处理器
	mqttHandler 	*mqtt.MqttClient
	//编解码处理器
	codecHandler 	*codec.Codec
	//运行环境配置
	cfg 			*config.EdgeConfig
}

func (c *AppCoreClient) Init() error {
	c.cfg = &config.EdgeConfig{}
	err := c.cfg.Load(c.appType)
	if err != nil {
		return errors.New("APP SDK init failed, err: " + err.Error())
	}
	c.codecHandler = codec.NewCodec(c.cfg.AppId, c.cfg.DeviceId, c.cfg.ThingId)
	clientId := fmt.Sprintf("%s/%s", c.cfg.DeviceId, c.cfg.AppId)
	url := fmt.Sprintf("%s://%s:%d", c.cfg.Protocol, c.cfg.HubAddr, c.cfg.HubPort)
	c.mqttHandler, err = mqtt.NewMqttClient(clientId, url, c.onConnectStatus)
	if err != nil {
		fmt.Println()
		//回滚已经初始化过的内容
		c.cfg = nil
		c.codecHandler = nil
		return errors.New("APP SDK init failed, err: " + err.Error())
	}
	return nil
}

func (c *AppCoreClient) Cleanup() {
	if c.mqttHandler != nil {
		c.mqttHandler.Stop()
		c.mqttHandler = nil
	}
	if c.codecHandler != nil {
		c.codecHandler = nil
	}
	if c.cfg != nil {
		c.cfg = nil
	}
}

func (c *AppCoreClient) Start() error {
	if c.mqttHandler == nil || c.codecHandler == nil || c.cfg == nil {
		return errors.New("APP SDK start failed, err: not init")
	}
	return c.mqttHandler.Start()
}

func (c *AppCoreClient) Stop() {
	if c.mqttHandler == nil || c.codecHandler == nil || c.cfg == nil {
		return
	}
	c.mqttHandler.Stop()
}

func (c *AppCoreClient) SendMessage(data *common.AppSdkMessageData) error {
	if c.mqttHandler == nil || c.codecHandler == nil || c.cfg == nil {
		return errors.New("APP SDK send message failed, err: not init")
	}
	if data == nil {
		return errors.New("APP SDK send message failed, err: invalid arguments")
	}
	var topicType string
	switch data.Type {
	case common.AppSdkMessageType_Property:
		topicType = codec.TopicType_PubProperty
	case common.AppSdkMessageType_Event:
		topicType = codec.TopicType_PubEvent
	case common.AppSdkMessageType_ServiceCall:
		topicType = codec.TopicType_PubService
	case common.AppSdkMessageType_ServiceReply:
		topicType = codec.TopicType_PubServiceReply
	default:
		return errors.New("APP SDK send message failed, err: unsupported message type")
	}
	var pubTopic string
	var pubData []byte
	if data.Type == common.AppSdkMessageType_Property || data.Type == common.AppSdkMessageType_Event ||
			data.Type == common.AppSdkMessageType_ServiceCall ||
			data.Type == common.AppSdkMessageType_ServiceReply {
		tempTopic, tempData, err := c.codecHandler.EncodeMessage(topicType, data.Payload)
		if err != nil {
			return err
		}
		pubTopic = tempTopic
		pubData = tempData
	}
	return c.mqttHandler.Publish(pubTopic, 0, pubData)
}

func (c *AppCoreClient) onConnectStatus(status bool, errMsg string) {
	if status {
		//Connected
		fmt.Println("APP SDK onConnectStatus called, status is connected")
		if c.mqttHandler == nil || c.codecHandler == nil || c.cfg == nil {
			fmt.Println("APP SDK onConnected subscribe topics failed, err: not init")
			return
		}
		//Callback connected event
		if c.eventCB != nil {
			evt := &common.AppSdkEventData{
				Type: common.EventType_Connected,
			}
			c.eventCB(evt, c.eventParam)
		}
		//Subscribe topics
		topics := make([]string, 0)
		tempTopic, err := c.codecHandler.EncodeTopic(codec.TopicType_SubProperty, "+")
		if err != nil {
			fmt.Printf("APP SDK onConnected EncodeTopic failed, topicType: %s, err: %s\n",
				codec.TopicType_SubProperty, err.Error())
		} else {
			topics = append(topics, tempTopic)
		}
		tempTopic, err = c.codecHandler.EncodeTopic(codec.TopicType_SubEvent, "+")
		if err != nil {
			fmt.Printf("APP SDK onConnected EncodeTopic failed, topicType: %s, err: %s\n",
				codec.TopicType_SubEvent, err.Error())
		} else {
			topics = append(topics, tempTopic)
		}
		for _, srvId := range c.serviceIds {
			tempTopic, err = c.codecHandler.EncodeTopic(codec.TopicType_SubService, srvId)
			if err != nil {
				fmt.Printf("APP SDK onConnected EncodeTopic failed, topicType: %s, err: %s\n",
					codec.TopicType_SubService, err.Error())
			} else {
				topics = append(topics, tempTopic)
			}
		}
		err = c.mqttHandler.SubscribeMultiple(topics, c.onRecvData)
		if err != nil {
			fmt.Println("APP SDK onConnected subscribe topics failed, err: " + err.Error())
		}
		fmt.Println("APP SDK onConnected subscribe topics success")
	} else {
		//Disconnected
		fmt.Println("APP SDK onConnectStatus called, status is disconnected, err: " + errMsg)
		if c.eventCB != nil {
			evt := &common.AppSdkEventData{
				Type: common.EventType_Disconnected,
			}
			c.eventCB(evt, c.eventParam)
		}
	}

}

func (c *AppCoreClient) onRecvData(topic string, payload []byte) {
	if c.mqttHandler == nil || c.codecHandler == nil || c.cfg == nil {
		fmt.Println("APP SDK onRecvData failed, err: not init")
		return
	}
	topicType, data, err := c.codecHandler.DecodeMessage(topic, payload)
	if err != nil {
		fmt.Println("APP SDK onRecvData DecodeMessage failed, err: " + err.Error())
		return
	}
	var msgType common.AppSdkMessageType
	switch topicType {
	case codec.TopicType_SubProperty:
		msgType = common.AppSdkMessageType_Property
	case codec.TopicType_SubEvent:
		msgType = common.AppSdkMessageType_Event
	case codec.TopicType_SubService:
		msgType = common.AppSdkMessageType_ServiceCall
	default:
		msgType = common.AppSdkMessageType_Unknown
	}
	msg := &common.AppSdkMessageData{
		Type: msgType,
		Payload: data,
	}
	if c.messageCB == nil {
		return
	}
	c.messageCB(msg, c.messageParam)
}