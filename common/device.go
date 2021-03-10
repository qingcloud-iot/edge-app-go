package common

//设备token状态
type TokenStatus string

const (
	//启用
	Enable  TokenStatus = "enabled"
	//禁止
	Disable TokenStatus = "disable"
)

/*
	边设备信息结构
*/
type EdgeLocalInfo struct {
	//应用
	AppId				string 			`json:"appId"`
	//边设备模型id
	ThingId 			string 			`json:"thingId"`
	//边设备id
	DeviceId 			string 			`json:"deviceId"`
}

/*
	子设备信息结构
*/
type EndpointInfo struct {
	//设备id
	DeviceId    		string			`json:"deviceId"`
	//设备名称
	DeviceName 			string			`json:"deviceName"`
	//模型id
	ThingId 			string			`json:"thingId"`
	//设备凭证
	Token 				string 			`json:"token"`
}