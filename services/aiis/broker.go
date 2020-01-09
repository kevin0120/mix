package aiis

import (
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/broker"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/utils"
)

const (
	SubjectResults     = "saturn.results.%s"
	SubjectResultsResp = "saturn.results.%s.response"
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

	s.dispatcherBus.Register(dispatcherbus.DispatcherNewTool, utils.CreateDispatchHandlerStruct(s.onNewTool))
	s.dispatcherBus.Register(dispatcherbus.DispatcherBrokerStatus, utils.CreateDispatchHandlerStruct(s.onBrokerStatus))
	return nil
}

func (s *BrokerClient) Stop() error {
	return nil
}

func (s *BrokerClient) Status() string {
	return s.broker.Status()
}

// 订阅工具结果
func (s *BrokerClient) collectTools() {
	for {
		select {
		case toolSN := <-s.toolCollector:
			err := s.broker.Subscribe(fmt.Sprintf(SubjectResultsResp, toolSN), s.onResultResp)
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
func (s *BrokerClient) onResultResp(message *broker.Message) ([]byte, error) {
	if message == nil {
		return nil, nil
	}

	payload := TransportPayload{}
	if err := json.Unmarshal(message.Body, &payload); err != nil {
		return nil, err
	}

	strData, _ := json.Marshal(payload.Data)
	rp := ResultPatch{}
	if err := json.Unmarshal(strData, &rp); err != nil {
		return nil, err
	}

	s.handlerResultPatch(rp)
	return nil, nil
}

func (s *BrokerClient) SendResult(result *PublishResult) error {
	str, _ := json.Marshal(TransportPayload{
		Method: TransportMethodResult,
		Data:   result,
	})

	return s.broker.Publish(fmt.Sprintf(SubjectResults, result.ToolSN), str)
}
