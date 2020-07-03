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
)

//消息回调
func onMessage(msg *common.AppSdkMessageData, param interface{}) {
	fmt.Println("onMessage called")
	if msg == nil {
		fmt.Println("onMessage failed, msg is nil")
		return
	}
	fmt.Println("msg type:", msg.Type)
	fmt.Println("msg payload:", msg.Payload)
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
			propMsg := &common.AppSdkMessageData{
				Type: common.AppSdkMessageType_Property,
			}
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
			propMsg.Payload, _ = json.Marshal(propData)
			err := cli.SendMessage(propMsg)
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
				evtMsg := &common.AppSdkMessageData{
					Type: common.AppSdkMessageType_Event,
				}
				evtMsg.Payload, _ = json.Marshal(evtData)
				err := cli.SendMessage(evtMsg)
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
	}
	client, err := edge_app_go.NewClient(options)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
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
	ctx, cancel := context.WithCancel(context.Background())
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