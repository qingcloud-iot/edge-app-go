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

func NewCodec(appId string, deviceId string, thingId string) *Codec {
	return &Codec{
		AppId: appId,
		DeviceId: deviceId,
		ThingId: thingId,
	}
}

type Codec struct {
	AppId 		string
	DeviceId    string
	ThingId 	string
}

//将SDK接口的消息数据编码成平台消息格式的数据
func (c *Codec) EncodeMessage(topicType string, payload []byte) (string, []byte, error) {
	switch topicType {
	case TopicType_SubProperty:
		data, err := c.encodePropertyMsg(payload)
		if err != nil {
			return "", nil, err
		}
		dstTopic := fmt.Sprintf(topicTemplate_SubProperty, c.AppId)
		return dstTopic, data, nil
	case TopicType_PubProperty:
		data, err := c.encodePropertyMsg(payload)
		if err != nil {
			return "", nil, err
		}
		dstTopic := fmt.Sprintf(topicTemplate_PubProperty, c.AppId)
		return dstTopic, data, nil
	case TopicType_SubEvent, TopicType_PubEvent:
		identifier, data, err := c.encodeEventMsg(payload)
		if err != nil {
			return "", nil, err
		}
		dstTopic, err := c.EncodeTopic(topicType, identifier)
		if err != nil {
			return "", nil, err
		}
		return dstTopic, data, nil
	case TopicType_PubService, TopicType_SubService:
		identifier, data, err := c.encodeServiceMsg(payload)
		if err != nil {
			return "", nil, err
		}
		dstTopic, err := c.EncodeTopic(topicType, identifier)
		if err != nil {
			return "", nil, err
		}
		return dstTopic, data, nil
	case TopicType_PubServiceReply:
		identifier, data, err := c.encodeServiceReplyMsg(payload)
		if err != nil {
			return "", nil, err
		}
		dstTopic, err := c.EncodeTopic(topicType, identifier)
		if err != nil {
			return "", nil, err
		}
		return dstTopic, data, nil
	}
	return "", nil, errors.New("Unsupported topic type: " + topicType)
}

//将平台消息格式的数据解码成SDK接口的消息数据
func (c *Codec) DecodeMessage(topic string, payload []byte) (string, []byte, error) {
	topicType, _, identifier, err := c.DecodeTopic(topic)
	if err != nil {
		return "", nil, err
	}
	switch topicType {
	case TopicType_SubProperty, TopicType_PubProperty:
		data, err := c.decodePropertyMsg(payload)
		if err != nil {
			return "", nil, err
		}
		return topicType, data, nil
	case TopicType_SubEvent, TopicType_PubEvent:
		data, err := c.decodeEventMsg(identifier, payload)
		if err != nil {
			return "", nil, err
		}
		return topicType, data, nil
	case TopicType_PubService, TopicType_SubService:
		data, err := c.decodeServiceMsg(identifier, payload)
		if err != nil {
			return "", nil, err
		}
		return topicType, data, nil
	case TopicType_PubServiceReply:
		data, err := c.decodeServiceReplyMsg(identifier, payload)
		if err != nil {
			return "", nil, err
		}
		return topicType, data, nil
	}
	return "", nil, errors.New("Unsupported topic type: " + topicType)
}

//编码Topic
func (c *Codec) EncodeTopic(topicType string, identifier string) (string, error) {
	if topicType == ""  {
		return "", errors.New("invalid arguments")
	}
	if (topicType == TopicType_SubEvent || topicType == TopicType_PubEvent || topicType == TopicType_PubService) &&
		identifier == "" {
		return "", errors.New("invalid arguments")
	}
	topic := ""
	switch topicType {
	case TopicType_SubProperty:
		topic = fmt.Sprintf(topicTemplate_SubProperty, c.AppId)
	case TopicType_PubProperty:
		topic = fmt.Sprintf(topicTemplate_PubProperty, c.AppId)
	case TopicType_SubEvent:
		topic = fmt.Sprintf(topicTemplate_SubEvent, c.AppId, identifier)
	case TopicType_PubEvent:
		topic = fmt.Sprintf(topicTemplate_PubEvent, c.AppId, identifier)
	case TopicType_PubService:
		topic = fmt.Sprintf(topicTemplate_PubService, c.AppId, identifier)
	case TopicType_SubService:
		topic = fmt.Sprintf(topicTemplate_SubService, c.ThingId, c.DeviceId, identifier)
	case TopicType_PubServiceReply:
		topic = fmt.Sprintf(topicTemplate_PubServiceReply, c.ThingId, c.DeviceId, identifier)
	default:
		return "", errors.New("unsupported topicType: " + topicType)
	}

	return topic, nil
}

//解码Topic
func (c *Codec) DecodeTopic(topic string) (string, string, string, error) {
	if topic == "" {
		return "", "", "", errors.New("invalid arguments")
	}
	dstUnits := strings.Split(topic, "/")
	if len(dstUnits) != 7 && len(dstUnits) != 8 {
		return "", "", "", errors.New("invalid topic format")
	}
	if dstUnits[0] != "" || (dstUnits[1] != "edge" && dstUnits[1] != "sys") {
		return "", "", "", errors.New("invalid topic format: " + topic)
	}
	//parse topic
	topicType := ""
	appId := dstUnits[2]
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
		return "", "", "", errors.New("invalid topic format: " + topic)
	}
	if (topicType == TopicType_SubEvent || topicType == TopicType_PubEvent || topicType == TopicType_PubService) &&
		identifier == "" {
		return "", "", "", errors.New("invalid identifier format: " + topic)
	}
	return topicType, appId, identifier, nil
}

func (c *Codec) encodePropertyMsg(payload []byte) ([]byte, error) {
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
		ModelId: c.ThingId,
		EntityId: c.DeviceId,
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

func (c *Codec) encodeEventMsg(payload []byte) (string, []byte, error) {
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
		ModelId: c.ThingId,
		EntityId: c.DeviceId,
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

func (c *Codec) encodeServiceMsg(payload []byte) (string, []byte, error) {
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
		ModelId: c.ThingId,
		EntityId: c.DeviceId,
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

func (c *Codec) decodePropertyMsg(payload []byte) ([]byte, error) {
	msg := &MdmpPropertyMsg{}
	err := json.Unmarshal(payload, msg)
	if err != nil {
		return nil, err
	}
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

func (c *Codec) decodeEventMsg(identifier string, payload []byte) ([]byte, error) {
	msg := &MdmpEventMsg{}
	err := json.Unmarshal(payload, msg)
	if err != nil {
		return nil, err
	}
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

func (c *Codec) decodeServiceMsg(identifier string, payload []byte) ([]byte, error) {
	msg := &MdmpServiceCallMsg{}
	err := json.Unmarshal(payload, msg)
	if err != nil {
		return nil, err
	}
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

func (c *Codec) decodeServiceReplyMsg(identifier string, payload []byte) ([]byte, error) {
	msg := &MdmpServiceReplyMsg{}
	err := json.Unmarshal(payload, msg)
	if err != nil {
		return nil, err
	}
	srv := &common.AppSdkMsgServiceCall{}
	srv.MessageId = msg.ID
	srv.Identifier = identifier
	srv.Params = msg.Data
	result, err := json.Marshal(srv)
	if err != nil {
		return nil, err
	}
	return result, nil
}