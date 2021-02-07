package meta

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/qingcloud-iot/edge-app-go/common"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	Metadata_Url_ChildDevice = "http://%s:%d/internal/data/childDevice"
)

func NewMetaClient(addr string, port int) *MetaClient {
	return &MetaClient{
		addr: addr,
		port: port,
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, netw string, addr string) (net.Conn, error) {
					c, err := net.DialTimeout(netw, addr, 5 * time.Second)
					if err != nil {
						return nil, err
					}
					deadline := time.Now().Add(5 * time.Second)
					c.SetDeadline(deadline)
					return c, nil
				},
			},
		},
	}
}

type MetaClient struct {
	addr 	string
	port 	int
	client  *http.Client
}

func (m *MetaClient) GetSubDevices() ([]*common.EndpointInfo, error) {
	url := fmt.Sprintf(Metadata_Url_ChildDevice, m.addr, m.port)
	resp, err := m.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	vTempInfos := make(map[string]interface{})
	err = json.Unmarshal([]byte(data), &vTempInfos)
	if err != nil {
		return nil, err
	}
	results := make([]*common.EndpointInfo, 0)
	for k, v := range vTempInfos {
		switch v.(type) {
		case string:
			k = strings.Replace(k, "/", "", 1)
			tempValue := v.(string)
			tempInfo := &common.EndpointInfo{}
			err = json.Unmarshal([]byte(tempValue), tempInfo)
			if err == nil {
				results = append(results, tempInfo)
			}
		default:
		}
	}
	return results, nil
}

