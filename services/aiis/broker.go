package aiis

import (
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/broker"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/utils"
)

const (
	SUBJECT_RESULTS      = "saturn.results.%s"
	SUBJECT_RESULTS_RESP = "saturn.results.%s.response"
)

func NewBrokerClient(d Diagnostic, bs IBrokerService, dp Dispatcher) *BrokerClient {
	return &BrokerClient{
		diag:          d,
		BaseTransport: BaseTransport{},
		toolCollector: make(chan string, 16),
		broker:        bs,
		dispatcherBus: dp,
	}
}

type BrokerClient struct {
	BaseTransport
	diag Diagnostic

	broker        IBrokerService
	dispatcherBus Dispatcher

	toolCollector chan string
	closing       chan struct{}
}

func (s *BrokerClient) Start() error {

	s.dispatcherBus.Register(dispatcherbus.DISPATCHER_NEW_TOOL, utils.CreateDispatchHandlerStruct(s.onNewTool))
	s.dispatcherBus.Register(dispatcherbus.DISPATCHER_BROKER_STATUS, utils.CreateDispatchHandlerStruct(s.onBrokerStatus))
	return nil
}

func (s *BrokerClient) Stop() error {
	return nil
}

func (s *BrokerClient) Status() string {
	// TODO: get broker conn status
	return TRANSPORT_STATUS_OFFLINE
}

// 订阅工具结果
func (s *BrokerClient) collectTools() {
	for {
		select {
		case toolSN := <-s.toolCollector:
			err := s.broker.Subscribe(fmt.Sprintf(SUBJECT_RESULTS_RESP, toolSN), s.onResultResp)
			if err != nil {
				s.diag.Error("Subscribe Failed", err)
			}

		case <-s.closing:
			return
		}
	}
}

// 接收Broker服务状态
func (s *BrokerClient) onBrokerStatus(data interface{}) {
	if data == nil {
		return
	}

	brokerStatus := data.(bool)
	if !brokerStatus {
		return
	}

	// Broker服务上线时，订阅工具结果
	go s.collectTools()
}

// 检测到新工具
func (s *BrokerClient) onNewTool(data interface{}) {
	if data == nil {
		return
	}

	toolSN := data.(string)

	// 加入待订阅列表
	s.toolCollector <- toolSN
}

// 收到上传结果反馈
func (s *BrokerClient) onResultResp(message *broker.BrokerMessage) ([]byte, error) {
	if message == nil {
		return nil, nil
	}

	payload := TransportPayload{}
	json.Unmarshal(message.Body, &payload)
	strData, _ := json.Marshal(payload.Data)
	rp := ResultPatch{}
	json.Unmarshal(strData, &rp)
	s.handlerResultPatch(rp)

	return nil, nil
}

func (c *BrokerClient) SendResult(result *AIISResult) error {
	str, _ := json.Marshal(TransportPayload{
		Method: TRANSPORT_METHOD_RESULT,
		Data:   result,
	})

	return c.broker.Publish(fmt.Sprintf(SUBJECT_RESULTS, result.ToolSN), str)
}
