package mqtt

import (
	"context"
	"errors"
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"time"
)

const (
	DefaultKeepAlive 		= 5 * time.Second
	DefaultWaitTimeout 		= 5 * time.Second
)

//第一次连接成功的回调（因为paho第一次连接上之后，内部有重连机制，所以只需要处理第一连接成功的重连，保证第一次能够连接成功）
type OnCollectedCallback func(bool, string)

//消息回调
type MessageCallback func(string, []byte)

func NewMqttClient(id string, url string, cb OnCollectedCallback) (*MqttClient, error) {
	mc := &MqttClient{}
	if url == "" || id == "" {
		return nil, errors.New("invalid arguments")
	}
	options := paho.NewClientOptions()
	options.AddBroker(url)
	options.SetClientID(id)
	options.SetCleanSession(true)
	options.SetAutoReconnect(true)
	options.SetKeepAlive(DefaultKeepAlive)
	options.SetOnConnectHandler(mc.onConnect)
	options.SetConnectionLostHandler(mc.onDisconnect)
	mc.client = paho.NewClient(options)
	mc.connectedCB = cb
	return mc, nil
}

type MqttClient struct {
	client 		paho.Client
	connectedCB OnCollectedCallback

	cancelCtx 	context.Context
	cancelFn 	context.CancelFunc
}

func (m *MqttClient) Start() error {
	m.cancelCtx, m.cancelFn = context.WithCancel(context.Background())
	go m.tryConnect()
	return nil
}

func (m *MqttClient) Stop() {
	if m.cancelFn != nil {
		m.cancelFn()
	}
	m.client = nil
}

func (m *MqttClient) Subscribe(topic string, qos int32, cb MessageCallback) error {
	if topic == "" || qos < 0 || qos > 2 || cb == nil {
		return errors.New("invalid arguments")
	}
	if token := m.client.Subscribe(topic, byte(qos), func(client paho.Client, message paho.Message) {
		cb(message.Topic(), message.Payload())
	}); token.WaitTimeout(DefaultWaitTimeout) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MqttClient) SubscribeMultiple(topics []string, cb MessageCallback) error {
	filters := make(map[string]byte)
	for _, topic := range topics {
		if topic == "" {
			continue
		}
		filters[topic] = 0
	}
	if token := m.client.SubscribeMultiple(filters, func(client paho.Client, message paho.Message) {
		cb(message.Topic(), message.Payload())
	}); token.WaitTimeout(DefaultWaitTimeout) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MqttClient) Unsubscribe(topics []string) error {
	if token := m.client.Unsubscribe(topics...); token.WaitTimeout(DefaultWaitTimeout) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MqttClient) Publish(topic string, qos int32, payload []byte) error {
	if topic == "" || qos < 0 || qos > 2 || payload == nil {
		return errors.New("invalid arguments")
	}
	if token := m.client.Publish(topic, byte(qos), false, payload); token.WaitTimeout(DefaultWaitTimeout) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MqttClient) tryConnect() {
	err := m.doConnect()
	if err == nil {
		return
	}
	//try to reconnect util connected
	tm := time.NewTimer(DefaultWaitTimeout)
	for {
		select {
		case <-m.cancelCtx.Done():
			return
		case <-tm.C:
			fmt.Println("try to connect edge_hub...")
			err = m.doConnect()
			if err == nil {
				return
			}
			tm.Reset(DefaultWaitTimeout)
		}
	}
}

func (m *MqttClient) doConnect() error {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MqttClient) onConnect(client paho.Client) {
	if m.connectedCB != nil {
		m.connectedCB(true, "")
	}
}

func (m *MqttClient) onDisconnect(client paho.Client, err error) {
	if m.connectedCB != nil {
		m.connectedCB(false, err.Error())
	}
}
