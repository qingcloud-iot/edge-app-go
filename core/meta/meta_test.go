package meta

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetaClient_GetSubDevices(t *testing.T) {
	assert := assert.New(t)
	metaClient := NewMetaClient("127.0.0.1", 9611)
	subDevices, err := metaClient.GetSubDevices()
	if !assert.Nil(err) {
		fmt.Println("test failed, err:", err.Error())
	} else {
		fmt.Println("test success, len:", len(subDevices))
	}
}