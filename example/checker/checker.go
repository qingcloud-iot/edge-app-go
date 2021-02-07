package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/qingcloud-iot/edge-app-go"
	"github.com/qingcloud-iot/edge-app-go/common"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	RANDOM_DATA_PROPERTY_ID = "random_data"
	RANDOM_DATA_EVENT_ID = "data_event"
	RANDOM_DATA_EVENT_PARAM_DATA = "event_value"
	RANDOM_DATA_SERVICE_CALL_ID = "test_app_call"
)

var client edge_app_go.Client

//消息回调
func onMessage(msg *common.AppSdkMessageData, param interface{}) {
	fmt.Println("onMessage called")
	if msg == nil {
		fmt.Println("onMessage failed, msg is nil")
		return
	}
	fmt.Println("msg type:", msg.Type)
	fmt.Println("msg thingId:", msg.ThingId)
	fmt.Println("msg deviceId:", msg.DeviceId)
	fmt.Println("msg payload:", string(msg.Payload))
	if msg.Type == common.AppSdkMessageType_ServiceCall {
		info := &common.AppSdkMsgServiceCall{}
		err := json.Unmarshal(msg.Payload, info)
		if err != nil {
			fmt.Println("onMessage Unmarshal failed,", err.Error())
			return
		}
		//只处理自己关心的服务调用
		if info.Identifier == RANDOM_DATA_SERVICE_CALL_ID {
			respData := &common.AppSdkMsgServiceReply{
				MessageId: info.MessageId,
				Identifier: info.Identifier,
				Code: 200,
				Params: info.Params,
			}
			respPayload, _ := json.Marshal(respData)
			err = client.SendMessage(common.AppSdkMessageType_ServiceReply, respPayload)
			if err != nil {
				fmt.Println("onMessage CallReply SendMessage failed,", err.Error())
				return
			}
			fmt.Println("onMessage CallReply SendMessage success")
		}
	}
}

//sdk事件回调
func onSdkEvent(evt *common.AppSdkEventData, param interface{}) {
	fmt.Println("onSdkEvent called")
	fmt.Println("evt type:", evt.Type)
}

//定时上报测试数据
func postData(ctx context.Context, cli edge_app_go.Client) {
	tm := time.NewTimer(3 * time.Second)
	for {
		select {
		case <- ctx.Done():
			return
		case <- tm.C:
			propData := make([]*common.AppSdkMsgProperty, 0)
			//获取100以内随机数
			rand.Seed(time.Now().Unix())
			value := rand.Intn(100)
			tempProp := &common.AppSdkMsgProperty{
				Identifier: RANDOM_DATA_PROPERTY_ID,
				Timestamp: time.Now().UnixNano() / 1e6,
				Value: strconv.Itoa(value),
			}
			propData = append(propData, tempProp)
			payload, _ := json.Marshal(propData)
			err := cli.SendMessage(common.AppSdkMessageType_Property, payload)
			if err != nil {
				fmt.Println("send property message failed, err:", err.Error())
			}
			if value >= 70 {
				//发送事件
				evtData := &common.AppSdkMsgEvent{
					Identifier: RANDOM_DATA_EVENT_ID,
					Timestamp: time.Now().UnixNano() / 1e6,
					Params: make(map[string]interface{}),
				}
				evtData.Params[RANDOM_DATA_EVENT_PARAM_DATA] = strconv.Itoa(value)
				evtPayload, _ := json.Marshal(evtData)
				err := cli.SendMessage(common.AppSdkMessageType_Event, evtPayload)
				if err != nil {
					fmt.Println("send event message failed, err:", err.Error())
				}
			}
			tm.Reset(3 * time.Second)
		}
	}

}

func main() {
	fmt.Println("checker start")
	options := &edge_app_go.Options{
		Type: common.AppSdkRuntimeType_Exec,
		MessageCB: onMessage,
		EventCB: onSdkEvent,
		ServiceIds: []string{RANDOM_DATA_SERVICE_CALL_ID},
	}
	//从环境变量读取子设备的模型id，此环境变量Key为开发者自定义
	endpointThingId := os.Getenv("ENDPOINT_THING_ID")
	if endpointThingId != "" {
		//订阅子设备的消息
		options.EndpointThingIds = []string{endpointThingId}
	}
	cli, err := edge_app_go.NewClient(options)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	client = cli
	err = client.Init()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = client.Start()
	if err != nil {
		fmt.Println(err.Error())
		client.Cleanup()
		return
	}
	//测试设置子设备的服务调用
	go func() {
		for {
			time.Sleep(5 * time.Second)
			req := &common.AppSdkMsgServiceCall{
				MessageId: "test_message_id",
				Identifier: "setTemperature",
				Params: make(map[string]interface{}),
			}
			req.Params["temperature"] = 35
			resp, err := client.CallEndpoint("iott-m45OOGLBYK", "iotd-3b6d0a2e-74db-4988-ba96-f8927b1c281a", req)
			if err != nil {
				fmt.Println("CallEndpoint failed, err:" + err.Error())
			} else {
				fmt.Println("CallEndpoint Response:", resp.Code)
			}
		}

	}()
	ctx, cancel := context.WithCancel(context.Background())
	//测试定时上报
	go postData(ctx, client)
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			goto Exit
		case syscall.SIGHUP:
		default:
			goto Exit
		}
	}
Exit:
	cancel()
	client.Cleanup()
	fmt.Println("checker exit")
}
