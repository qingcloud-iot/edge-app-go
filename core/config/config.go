package config

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/qingcloud-iot/edge-app-go/common"
	"io/ioutil"
	"os"
	"strconv"
)

//二进制应用加载的配置
var edgeconf = flag.String("edgeconfig", "", "edge app config path")

//Docker应用加载环境参数
const (
	//EdgeHub协议类型
	ENV_EDGE_HUB_PROTOCOL 	= "EDGE_HUB_PROTO"
	//EdgeHub地址
	ENV_EDGE_HUB_HOST  		= "EDGE_HUB_HOST"
	//EdgeHub端口
	ENV_EDGE_HUB_PORT  		= "EDGE_HUB_PORT"
	//应用id
	ENV_EDGE_APP_ID    		= "EDGE_APP_ID"
	//边设备id
	ENV_EDGE_DEVICE_ID 		= "EDGE_DEVICE_ID"
	//设备模型id
	ENV_EDGE_THING_ID  		= "EDGE_THING_ID"
	//消息模式
	ENV_EDGE_MSG_MODE 		= "EDGE_MSG_MODE"
)

//二进制应用初始化配置文件结构
type EdgeConfig struct {
	//EdgeHub协议类型
	Protocol string 		`json:"protocol"`
	//EdgeHub地址
	HubAddr  string 		`json:"hubAddr"`
	//EdgeHub端口
	HubPort  int    		`json:"hubPort"`
	//应用id
	AppId    string 		`json:"appId"`
	//设备id
	DeviceId string 		`json:"deviceId"`
	//设备模型id
	ThingId  string 		`json:"thingId"`
}

func (c *EdgeConfig) Load(appType common.AppSdkRuntimeType) error {
	if appType == common.AppSdkRuntimeType_Exec {
		//加载运行参数
		flag.Parse()
		//加载配置
		jsonFile, err := ioutil.ReadFile(*edgeconf)
		if err != nil {
			return err
		}
		err = json.Unmarshal(jsonFile, c)
		if err != nil {
			return err
		}

	} else if appType == common.AppSdkRuntimeType_Docker {
		//加载环境变量
		c.Protocol = os.Getenv(ENV_EDGE_HUB_PROTOCOL)
		if c.Protocol == "" {
			c.Protocol = "tcp"
		}
		c.HubAddr = os.Getenv(ENV_EDGE_HUB_HOST)
		portStr := os.Getenv(ENV_EDGE_HUB_PORT)
		if portStr == "" {
			c.HubPort = 1883
		} else {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				port = 1883
			}
			c.HubPort = port
		}
		c.AppId = os.Getenv(ENV_EDGE_APP_ID)
		c.DeviceId = os.Getenv(ENV_EDGE_DEVICE_ID)
		c.ThingId = os.Getenv(ENV_EDGE_THING_ID)
	} else {
		return errors.New("Application type is not supported, appType: " + strconv.Itoa(int(appType)))
	}
	if c.HubAddr == "" || c.AppId == "" || c.DeviceId == "" {
		return errors.New("Load config failed, config param should not be empty")
	}
	return nil
}