package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qingcloud-iot/edge-app-go/common"
	"github.com/qingcloud-iot/edge-app-go/core/codec"
	"github.com/qingcloud-iot/edge-app-go/core/config"
	"github.com/qingcloud-iot/edge-app-go/core/meta"
	"github.com/qingcloud-iot/edge-app-go/core/mqtt"
	"github.com/satori/go.uuid"
	"time"
)

func NewAppCoreClient(appType common.AppSdkRuntimeType, msgCB common.AppSdkMessageCB, msgParam interface{},
						evtCB common.AppSdkEventCB, evtParam interface{}, srvIds []string, thingIds []string) *AppCoreClient {
	return &AppCoreClient{
		appType: 		appType,
		messageCB: 		msgCB,
		messageParam: 	msgParam,
		eventCB: 		evtCB,
		eventParam:  	evtParam,
		serviceIds: 	srvIds,
		epThingIds: 	thingIds,
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
	//子设备模型id数组
	epThingIds 		[]string
	//mqtt协议处理器
	mqttHandler 	*mqtt.MqttClient
	//编解码处理器
	codecHandler 	*codec.Codec
	//运行环境配置
	cfg 			*config.EdgeConfig
	//metadata访问客户端
	metaHandler 	*meta.MetaClient
}

func (c *AppCoreClient) Init() error {
	c.cfg = &config.EdgeConfig{}
	err := c.cfg.Load(c.appType)
	if err != nil {
		return errors.New("APP SDK init failed, err: " + err.Error())
	}
	c.codecHandler = codec.NewCodec(c.cfg.AppId, c.cfg.DeviceId, c.cfg.ThingId, c.cfg.ProxyMode)
	clientId := fmt.Sprintf("%s/%s", c.cfg.DeviceId, c.cfg.AppId)
	url := fmt.Sprintf("%s://%s:%d", c.cfg.Protocol, c.cfg.HubAddr, c.cfg.HubPort)
	c.mqttHandler, err = mqtt.NewMqttClient(clientId, url, c.onConnectStatus)
	c.metaHandler = meta.NewMetaClient(c.cfg.HubAddr, 9611)
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

func (c *AppCoreClient) SendMessage(msgType common.AppSdkMessageType, payload []byte) error {
	if c.mqttHandler == nil || c.codecHandler == nil || c.cfg == nil {
		return errors.New("APP SDK send message failed, err: not init")
	}
	if payload == nil {
		return errors.New("APP SDK send message failed, err: invalid arguments")
	}
	var topicType string
	switch msgType {
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
	if msgType == common.AppSdkMessageType_Property || msgType == common.AppSdkMessageType_Event ||
		msgType == common.AppSdkMessageType_ServiceCall ||
		msgType == common.AppSdkMessageType_ServiceReply {
		tempTopic, tempData, err := c.codecHandler.EncodeMessage(topicType, c.cfg.ThingId, c.cfg.DeviceId, payload)
		if err != nil {
			return err
		}
		pubTopic = tempTopic
		pubData = tempData
	}
	return c.mqttHandler.Publish(pubTopic, 0, pubData)
}

func (c *AppCoreClient) GetEdgeDeviceInfo() (*common.EdgeLocalInfo, error) {
	info := &common.EdgeLocalInfo{}
	info.AppId = c.cfg.AppId
	info.ThingId = c.cfg.ThingId
	info.DeviceId = c.cfg.DeviceId
	return info, nil
}

func (c *AppCoreClient) GetEndpointInfos() ([]*common.EndpointInfo, error) {
	if c.metaHandler == nil {
		return nil, errors.New("APP SDK send message failed, err: not init")
	}
	return c.metaHandler.GetSubDevices()
}

func (c *AppCoreClient) CallEndpoint(thingId string, deviceId string, req *common.AppSdkMsgServiceCall) (*common.AppSdkMsgServiceReply, error) {
	if thingId == "" || deviceId == "" || req == nil || req.Identifier == "" {
		return nil, errors.New("APP SDK CallEndpoint failed, err: invalid arguments")
	}
	if req.MessageId == "" {
		req.MessageId = uuid.NewV1().String()
	}
	//encode message
	tempData, _ := json.Marshal(req)
	callTopic, callPayload, err := c.codecHandler.EncodeMessage(codec.TopicType_PubService, thingId, deviceId, tempData)
	if err != nil {
		return nil, errors.New("APP SDK CallEndpoint failed, err: " + err.Error())
	}
	//generate reply topic
	replyTopic, err := c.codecHandler.EncodeTopic(codec.TopicType_SubServiceReply, req.Identifier, thingId, deviceId)
	if err != nil {
		return nil, errors.New("APP SDK CallEndpoint failed, err: " + err.Error())
	}
	exitCh := make(chan error)
	replyCh := make(chan *common.AppSdkMsgServiceReply)
	err = c.mqttHandler.Subscribe(replyTopic, 0, func(topic string, payload []byte) {
		if c.mqttHandler == nil || c.codecHandler == nil || c.cfg == nil {
			exitCh <- errors.New("APP SDK CallEndpoint callback failed, err: not init")
			return
		}
		topicType, _, _, data, err := c.codecHandler.DecodeMessage(topic, payload)
		if err != nil {
			exitCh <- errors.New("APP SDK CallEndpoint callback failed, err: " + err.Error())
			return
		}
		if topicType != codec.TopicType_SubServiceReply {
			exitCh <- errors.New("APP SDK CallEndpoint callback failed, err: not service reply")
			return
		}
		replyMsg := &common.AppSdkMsgServiceReply{}
		err = json.Unmarshal(data, replyMsg)
		if err != nil {
			exitCh <- errors.New("APP SDK CallEndpoint callback failed, err: " + err.Error())
			return
		}
		if replyMsg.MessageId != req.MessageId {
			//Match reply message id failed
			return
		}
		replyCh <- replyMsg
	})
	if err != nil {
		return nil, errors.New("APP SDK CallEndpoint failed, err: " + err.Error())
	}
	defer c.mqttHandler.Unsubscribe([]string{replyTopic})
	err = c.mqttHandler.Publish(callTopic, 0, callPayload)
	if err != nil {
		return nil, errors.New("APP SDK CallEndpoint failed, err: " + err.Error())
	}
	select {
	case value := <-replyCh:
		return value, nil
	case err := <- exitCh:
		return nil, err
	case <- time.After(time.Second * 5):
		return nil, errors.New("APP SDK CallEndpoint failed, err: call timeout")
	}
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
		//Subscribe topics of edge device
		topics := make([]string, 0)
		tempTopic, err := c.codecHandler.EncodeTopic(codec.TopicType_SubProperty, "+", c.cfg.ThingId, c.cfg.DeviceId)
		if err != nil {
			fmt.Printf("APP SDK onConnected EncodeTopic failed, topicType: %s, err: %s\n",
				codec.TopicType_SubProperty, err.Error())
		} else {
			topics = append(topics, tempTopic)
		}
		tempTopic, err = c.codecHandler.EncodeTopic(codec.TopicType_SubEvent, "+", c.cfg.ThingId, c.cfg.DeviceId)
		if err != nil {
			fmt.Printf("APP SDK onConnected EncodeTopic failed, topicType: %s, err: %s\n",
				codec.TopicType_SubEvent, err.Error())
		} else {
			topics = append(topics, tempTopic)
		}
		for _, srvId := range c.serviceIds {
			tempTopic, err = c.codecHandler.EncodeTopic(codec.TopicType_SubService, srvId, c.cfg.ThingId, c.cfg.DeviceId)
			if err != nil {
				fmt.Printf("APP SDK onConnected EncodeTopic failed, topicType: %s, err: %s\n",
					codec.TopicType_SubService, err.Error())
			} else {
				topics = append(topics, tempTopic)
			}
		}
		//非代理模式下，可以直接订阅子设备的模型消息
		if !c.cfg.ProxyMode {
			//Subscribe topics of endpoints
			for _, thingId := range c.epThingIds {
				tempTopic, err := c.codecHandler.EncodeTopic(codec.TopicType_SubProperty, "+", thingId, "+")
				if err != nil {
					fmt.Printf("APP SDK onConnected EncodeTopic for endpoints failed, topicType: %s, err: %s\n",
						codec.TopicType_SubProperty, err.Error())
				} else {
					topics = append(topics, tempTopic)
				}
				tempTopic, err = c.codecHandler.EncodeTopic(codec.TopicType_SubEvent, "+", thingId, "+")
				if err != nil {
					fmt.Printf("APP SDK onConnected EncodeTopic for endpoints failed, topicType: %s, err: %s\n",
						codec.TopicType_SubEvent, err.Error())
				} else {
					topics = append(topics, tempTopic)
				}
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
	topicType, thingId, deviceId, data, err := c.codecHandler.DecodeMessage(topic, payload)
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
		ThingId: thingId,
		DeviceId: deviceId,
		Payload: data,
	}
	if c.messageCB == nil {
		return
	}
	c.messageCB(msg, c.messageParam)
}