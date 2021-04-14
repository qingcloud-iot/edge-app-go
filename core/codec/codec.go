package codec

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qingcloud-iot/edge-app-go/common"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

func NewCodec(appId string, deviceId string, thingId string, proxyMode bool) *Codec {
	return &Codec{
		AppId: appId,
		DeviceId: deviceId,
		ThingId: thingId,
		ProxyMode: proxyMode,
	}
}

type Codec struct {
	AppId 		string
	DeviceId    string
	ThingId 	string
	ProxyMode   bool
}

//将SDK接口的消息数据编码成平台消息格式的数据
func (c *Codec) EncodeMessage(topicType string, thingId string, deviceId string, payload []byte) (string, []byte, error) {
	switch topicType {
	case TopicType_SubProperty:
		data, err := c.encodePropertyMsg(thingId, deviceId, payload)
		if err != nil {
			return "", nil, err
		}
		dstTopic, err := c.EncodeTopic(topicType, "", thingId, deviceId, c.ProxyMode)
		if err != nil {
			return "", nil, err
		}
		return dstTopic, data, nil
	case TopicType_PubProperty:
		data, err := c.encodePropertyMsg(thingId, deviceId, payload)
		if err != nil {
			return "", nil, err
		}
		dstTopic, err := c.EncodeTopic(topicType, "", thingId, deviceId, c.ProxyMode)
		if err != nil {
			return "", nil, err
		}
		return dstTopic, data, nil
	case TopicType_SubEvent, TopicType_PubEvent:
		identifier, data, err := c.encodeEventMsg(thingId, deviceId, payload)
		if err != nil {
			return "", nil, err
		}
		dstTopic, err := c.EncodeTopic(topicType, identifier, thingId, deviceId, c.ProxyMode)
		if err != nil {
			return "", nil, err
		}
		return dstTopic, data, nil
	case TopicType_PubService, TopicType_SubService:
		identifier, data, err := c.encodeServiceMsg(thingId, deviceId, payload)
		if err != nil {
			return "", nil, err
		}
		dstTopic, err := c.EncodeTopic(topicType, identifier, thingId, deviceId, c.ProxyMode)
		if err != nil {
			return "", nil, err
		}
		return dstTopic, data, nil
	case TopicType_PubServiceReply:
		identifier, data, err := c.encodeServiceReplyMsg(payload)
		if err != nil {
			return "", nil, err
		}
		dstTopic, err := c.EncodeTopic(topicType, identifier, thingId, deviceId, c.ProxyMode)
		if err != nil {
			return "", nil, err
		}
		return dstTopic, data, nil
	}
	return "", nil, errors.New("Unsupported topic type: " + topicType)
}

//将平台消息格式的数据解码成SDK接口的消息数据
func (c *Codec) DecodeMessage(topic string, payload []byte) (string, string, string, []byte, error) {
	topicType, thingId, deviceId, identifier, err := c.DecodeTopic(topic, c.ProxyMode)
	if err != nil {
		return "", "", "", nil, err
	}
	switch topicType {
	case TopicType_SubProperty, TopicType_PubProperty:
		msg, err := c.decodePropertyMsg(payload)
		if err != nil {
			return "", "", "", nil, err
		}
		tempDeviceId := c.getDeviceIdFromModel(msg.Metadata)
		if tempDeviceId != "" {
			deviceId = tempDeviceId
		}
		data, err := c.encodePropertySdkData(msg)
		if err != nil {
			return "", "", "", nil, err
		}
		return topicType, thingId, deviceId, data, nil
	case TopicType_SubEvent, TopicType_PubEvent:
		msg, err := c.decodeEventMsg(payload)
		if err != nil {
			return "", "", "", nil, err
		}
		tempDeviceId := c.getDeviceIdFromModel(msg.Metadata)
		if tempDeviceId != "" {
			deviceId = tempDeviceId
		}
		data, err := c.encodeEventSdkData(identifier, msg)
		if err != nil {
			return "", "", "", nil, err
		}
		return topicType, thingId, deviceId, data, nil
	case TopicType_PubService, TopicType_SubService:
		msg, err := c.decodeServiceCallMsg(payload)
		if err != nil {
			return "", "", "", nil, err
		}
		tempDeviceId := c.getDeviceIdFromModel(msg.Metadata)
		if tempDeviceId != "" {
			deviceId = tempDeviceId
		}
		data, err := c.encodeServiceCallSdkData(identifier, msg)
		if err != nil {
			return "", "", "", nil, err
		}
		return topicType, thingId, deviceId, data, nil
	case TopicType_PubServiceReply, TopicType_SubServiceReply:
		msg, err := c.decodeServiceReplyMsg(payload)
		if err != nil {
			return "", "", "", nil, err
		}
		data, err := c.encodeServiceReplySdkData(identifier, msg)
		if err != nil {
			return "", "", "", nil, err
		}
		return topicType, thingId, deviceId, data, nil
	}
	return "", "", "", nil, errors.New("Unsupported topic type: " + topicType)
}

//编码Topic
func (c *Codec) EncodeTopic(topicType string, identifier string, thingId string, deviceId string, proxyMode bool) (string, error) {
	if topicType == ""  {
		return "", errors.New("invalid arguments")
	}
	topic := ""
	if proxyMode {
		//代理模式
		switch topicType {
		case TopicType_SubProperty:
			topic = fmt.Sprintf(topicTemplateV1_SubProperty, c.AppId)
		case TopicType_PubProperty:
			topic = fmt.Sprintf(topicTemplateV1_PubProperty, c.AppId)
		case TopicType_SubEvent:
			topic = fmt.Sprintf(topicTemplateV1_SubEvent, c.AppId, identifier)
		case TopicType_PubEvent:
			topic = fmt.Sprintf(topicTemplateV1_PubEvent, c.AppId, identifier)
		case TopicType_PubService:
			topic = fmt.Sprintf(topicTemplateV1_PubService, c.AppId, identifier)
		case TopicType_SubService:
			topic = fmt.Sprintf(topicTemplateV1_SubService, thingId, deviceId, identifier)
		case TopicType_PubServiceReply:
			topic = fmt.Sprintf(topicTemplateV1_PubServiceReply, thingId, deviceId, identifier)
		case TopicType_SubServiceReply:
			topic = fmt.Sprintf(topicTemplateV1_SubServiceReply, thingId, deviceId, identifier)
		default:
			return "", errors.New("unsupported topicType: " + topicType)
		}
	} else {
		//普通模式
		switch topicType {
		case TopicType_SubProperty:
			topic = fmt.Sprintf(topicTemplateV2_SubProperty, thingId, deviceId)
		case TopicType_PubProperty:
			topic = fmt.Sprintf(topicTemplateV2_PubProperty, thingId, deviceId)
		case TopicType_SubEvent:
			topic = fmt.Sprintf(topicTemplateV2_SubEvent, thingId, deviceId, identifier)
		case TopicType_PubEvent:
			topic = fmt.Sprintf(topicTemplateV2_PubEvent, thingId, deviceId, identifier)
		case TopicType_PubService:
			topic = fmt.Sprintf(topicTemplateV2_PubService, thingId, deviceId, identifier)
		case TopicType_SubService:
			topic = fmt.Sprintf(topicTemplateV2_SubService, thingId, deviceId, identifier)
		case TopicType_PubServiceReply:
			topic = fmt.Sprintf(topicTemplateV2_PubServiceReply, thingId, deviceId, identifier)
		case TopicType_SubServiceReply:
			topic = fmt.Sprintf(topicTemplateV2_SubServiceReply, thingId, deviceId, identifier)
		default:
			return "", errors.New("unsupported topicType: " + topicType)
		}

	}
	return topic, nil
}

//解码Topic
func (c *Codec) DecodeTopic(topic string, proxyMode bool) (string, string, string, string, error) {
	if topic == "" {
		return "", "", "", "", errors.New("invalid arguments")
	}
	if proxyMode {
		//代理模式
		dstUnits := strings.Split(topic, "/")
		if len(dstUnits) != 7 && len(dstUnits) != 8 {
			return "", "", "", "", errors.New("invalid topic format")
		}
		if dstUnits[0] != "" || (dstUnits[1] != "edge" && dstUnits[1] != "sys") {
			return "", "", "", "", errors.New("invalid topic format: " + topic)
		}
		//parse topic
		topicType := ""
		identifier := ""
		if dstUnits[1] == "edge" {
			if dstUnits[4] == "property" {
				if dstUnits[6] == "post" {
					topicType = TopicType_SubProperty
				} else if dstUnits[6] == "control" {
					topicType = TopicType_PubProperty
				}
			} else if dstUnits[4] == "event" {
				if dstUnits[6] == "post" {
					topicType = TopicType_SubEvent
				} else if dstUnits[6] == "control" {
					topicType = TopicType_PubEvent
				}
				identifier = dstUnits[5]
			} else if dstUnits[4] == "service" {
				if dstUnits[6] == "call" {
					topicType = TopicType_PubService
				}
				identifier = dstUnits[5]
			}
		} else {
			if dstUnits[5] == "service" {
				if dstUnits[7] == "call" {
					topicType = TopicType_SubService
				} else if dstUnits[7] == "call_reply" {
					topicType = TopicType_PubServiceReply
				}
				identifier = dstUnits[6]
			}
		}
		if topicType == "" {
			return "", "", "", "", errors.New("invalid topic format: " + topic)
		}
		return topicType, c.ThingId, c.DeviceId, identifier, nil
	} else {
		//普通模式
		dstUnits := strings.Split(topic, "/")
		if len(dstUnits) != 8 {
			return "", "", "", "", errors.New("invalid topic format")
		}
		if dstUnits[0] != "" || dstUnits[1] != "sys" {
			return "", "", "", "", errors.New("invalid topic format: " + topic)
		}
		//parse topic
		topicType := ""
		thingId := dstUnits[2]
		deviceId := dstUnits[3]
		identifier := ""
		if dstUnits[5] == "property" {
			if dstUnits[7] == "post" {
				topicType = TopicType_SubProperty
			}
		} else if dstUnits[5] == "event" {
			if dstUnits[7] == "post" {
				topicType = TopicType_SubEvent
			}
			identifier = dstUnits[6]
		} else if dstUnits[5] == "service" {
			if dstUnits[7] == "call" {
				topicType = TopicType_SubService
			} else if dstUnits[7] == "call_reply" {
				topicType = TopicType_SubServiceReply
			}
			identifier = dstUnits[6]
		}
		if topicType == "" {
			return "", "", "", "", errors.New("invalid topic format: " + topic)
		}
		return topicType, thingId, deviceId, identifier, nil
	}
}

func (c *Codec) encodePropertyMsg(thingId string, deviceId string, payload []byte) ([]byte, error) {
	props := make([]*common.AppSdkMsgProperty, 0)
	err := json.Unmarshal(payload, &props)
	if err != nil {
		return nil, err
	}
	if len(props) == 0 {
		return nil, errors.New("properties is empty")
	}
	now := time.Now().UnixNano() / 1e6
	msg := &MdmpPropertyMsg{}
	msg.ID = uuid.NewV1().String()
	msg.Version = DefaultMessageVersion
	msg.Type = MessageTypeTemplate_Property
	msg.Metadata = &ModelMetadata{
		ModelId: thingId,
		EntityId: deviceId,
		Source: make([]string, 0),
		EpochTime: now,
	}
	msg.Params = make(map[string]*ModelPropertyData)
	for _, prop := range props {
		tempData := &ModelPropertyData{}
		tempData.Value = prop.Value
		tempData.Time = prop.Timestamp
		msg.Params[prop.Identifier] = tempData
	}
	result, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Codec) encodeEventMsg(thingId string, deviceId string, payload []byte) (string, []byte, error) {
	evt := &common.AppSdkMsgEvent{}
	err := json.Unmarshal(payload, evt)
	if err != nil {
		return "", nil, err
	}
	now := time.Now().UnixNano() / 1e6
	msg := &MdmpEventMsg{}
	msg.ID = uuid.NewV1().String()
	msg.Version = DefaultMessageVersion
	msg.Type = fmt.Sprintf(MessageTypeTemplate_Event, evt.Identifier)
	msg.Metadata = &ModelMetadata{
		ModelId: thingId,
		EntityId: deviceId,
		Source: make([]string, 0),
		EpochTime: now,
	}
	msg.Params = &ModelEventData{
		Time: evt.Timestamp,
		Value: evt.Params,
	}
	result, err := json.Marshal(msg)
	if err != nil {
		return "", nil, err
	}
	return evt.Identifier, result, nil
}

func (c *Codec) encodeServiceMsg(thingId string, deviceId string,payload []byte) (string, []byte, error) {
	srv := &common.AppSdkMsgServiceCall{}
	err := json.Unmarshal(payload, srv)
	if err != nil {
		return "", nil, err
	}
	msg := &MdmpServiceCallMsg{}
	msg.ID = srv.MessageId
	msg.Version = DefaultMessageVersion
	msg.Type = fmt.Sprintf(MessageTypeTemplate_Service, srv.Identifier)
	msg.Metadata = &ServiceMetadata{
		ModelId: thingId,
		EntityId: deviceId,
	}
	msg.Params = srv.Params
	result, err := json.Marshal(msg)
	if err != nil {
		return "", nil, err
	}
	return srv.Identifier, result, nil
}

func (c *Codec) encodeServiceReplyMsg(payload []byte) (string, []byte, error) {
	reply := &common.AppSdkMsgServiceReply{}
	err := json.Unmarshal(payload, reply)
	if err != nil {
		return "", nil, err
	}
	msg := &MdmpServiceReplyMsg{}
	msg.ID = reply.MessageId
	msg.Version = DefaultMessageVersion
	msg.Code = reply.Code
	msg.Data = reply.Params
	result, err := json.Marshal(msg)
	if err != nil {
		return "", nil, err
	}
	return reply.Identifier, result, nil
}

func (c *Codec) decodePropertyMsg(payload []byte) (*MdmpPropertyMsg, error) {
	msg := &MdmpPropertyMsg{}
	err := json.Unmarshal(payload, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Codec) encodePropertySdkData(msg *MdmpPropertyMsg) ([]byte, error) {
	props := make([]*common.AppSdkMsgProperty, 0)
	for k, v := range msg.Params {
		tempProp := &common.AppSdkMsgProperty{}
		tempProp.Identifier = k
		tempProp.Value = v.Value
		tempProp.Timestamp = v.Time
		props = append(props, tempProp)
	}
	result, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Codec) decodeEventMsg(payload []byte) (*MdmpEventMsg, error) {
	msg := &MdmpEventMsg{}
	err := json.Unmarshal(payload, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Codec) encodeEventSdkData(identifier string, msg *MdmpEventMsg) ([]byte, error) {
	evt := &common.AppSdkMsgEvent{}
	evt.Identifier = identifier
	evt.Timestamp = msg.Params.Time
	evt.Params = msg.Params.Value
	result, err := json.Marshal(evt)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Codec) decodeServiceCallMsg(payload []byte) (*MdmpServiceCallMsg, error) {
	msg := &MdmpServiceCallMsg{}
	err := json.Unmarshal(payload, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Codec) encodeServiceCallSdkData(identifier string, msg *MdmpServiceCallMsg) ([]byte, error) {
	srv := &common.AppSdkMsgServiceCall{}
	srv.MessageId = msg.ID
	srv.Identifier = identifier
	srv.Params = msg.Params
	result, err := json.Marshal(srv)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Codec) decodeServiceReplyMsg(payload []byte) (*MdmpServiceReplyMsg, error) {
	msg := &MdmpServiceReplyMsg{}
	err := json.Unmarshal(payload, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Codec) encodeServiceReplySdkData(identifier string, msg *MdmpServiceReplyMsg) ([]byte, error) {
	srv := &common.AppSdkMsgServiceReply{}
	srv.MessageId = msg.ID
	srv.Identifier = identifier
	srv.Code = msg.Code
	srv.Params = msg.Data
	result, err := json.Marshal(srv)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Codec) getDeviceIdFromModel(md interface{}) string {
	if v, ok := md.(map[string]interface{}); ok {
		return v["entityId"].(string)
	}
	return ""
}