package common

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalMsgProperty(t *testing.T) {
	assert := assert.New(t)
	expectedResult := `[{"identifier":"id_prop_01","timestamp":1593274999806,"value":"aaaaaa"},{"identifier":"id_prop_02","timestamp":1593274999806,"value":"bbbbbb"}]`
	props := make([]*AppSdkMsgProperty, 0)
	prop01 := &AppSdkMsgProperty{
		Identifier: "id_prop_01",
		Timestamp: 1593274999806,
		Value: "aaaaaa",
	}
	prop02 := &AppSdkMsgProperty{
		Identifier: "id_prop_02",
		Timestamp: 1593274999806,
		Value: "bbbbbb",
	}
	props = append(props, prop01, prop02)
	payload, err := json.Marshal(props)
	if !assert.Nil(err) {
		return
	}
	fmt.Println("test marshal props:", string(payload))
	assert.Equal(expectedResult, string(payload))
}

func TestUnmarshalMsgProperty(t *testing.T) {
	assert := assert.New(t)
	src := `[{"identifier":"id_prop_01","timestamp":1593274999806,"value":"aaaaaa"},{"identifier":"id_prop_02","timestamp":1593274999806,"value":"bbbbbb"}]`
	props := make([]*AppSdkMsgProperty, 0)
	err := json.Unmarshal([]byte(src), &props)
	if !assert.Nil(err) {
		return
	}
	fmt.Println("test unmarshal prop success.")
}

func TestMarshalMsgEvent(t *testing.T) {
	assert := assert.New(t)
	expectedResult := `{"identifier":"test_event_001","timestamp":1593274999806,"params":{"param1":"aaa","param2":20,"param3":"ccc"}}`
	evt := &AppSdkMsgEvent{
		Identifier: "test_event_001",
		Timestamp: 1593274999806,
		Params: map[string]interface{}{
			"param1": "aaa",
			"param2": 20,
			"param3": "ccc",
		},
	}
	payload, err := json.Marshal(evt)
	if !assert.Nil(err) {
		return
	}
	fmt.Println("test marshal event:", string(payload))
	assert.Equal(expectedResult, string(payload))
}

func TestUnmarshalMsgEvent(t *testing.T) {
	assert := assert.New(t)
	src := `{"identifier":"test_event_001","timestamp":1593274999806,"params":{"param1":"aaa","param2":20,"param3":"ccc"}}`
	evt := &AppSdkMsgEvent{}
	err := json.Unmarshal([]byte(src), evt)
	if !assert.Nil(err) {
		return
	}
	fmt.Println("test unmarshal event success.")
}

func TestMarshalMsgServiceCall(t *testing.T) {
	assert := assert.New(t)
	expectedResult := `{"identifier":"test_service_001","params":{"param1":"aaa","param2":20,"param3":"ccc"}}`
	evt := &AppSdkMsgServiceCall{
		Identifier: "test_service_001",
		Params: map[string]interface{}{
			"param1": "aaa",
			"param2": 20,
			"param3": "ccc",
		},
	}
	payload, err := json.Marshal(evt)
	if !assert.Nil(err) {
		return
	}
	fmt.Println("test marshal service call:", string(payload))
	assert.Equal(expectedResult, string(payload))
}

func TestUnmarshalMsgServiceCall(t *testing.T) {
	assert := assert.New(t)
	src := `{"identifier":"test_service_001","params":{"param1":"aaa","param2":20,"param3":"ccc"}}`
	srvcall := &AppSdkMsgServiceCall{}
	err := json.Unmarshal([]byte(src), srvcall)
	if !assert.Nil(err) {
		return
	}
	fmt.Println("test unmarshal service call success.")
}