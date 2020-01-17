package rush

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/transport"
)

type ResultHandler func(result *aiis.PublishResult)

type BaseTransport struct {
	aiis.BaseTransport
}

// 发送结果上传反馈
func (s *BaseTransport) SendResultPatch(toolSN string, rp *aiis.ResultPatch) error {
	trans := s.ITransportService
	if trans == nil {
		return errors.New("trans Is Empty")
	}
	data, _ := json.Marshal(aiis.TransportPayload{
		Method: aiis.TransportMethodResultPatch,
		Data:   rp,
	})

	return trans.SendMessage(fmt.Sprintf(SUBJECT_RESULTS_RESP, toolSN), data)
}

// 发送服务状态
func (s *BaseTransport) SendServiceStatus(status *aiis.ServiceStatus) error {
	trans := s.ITransportService
	if trans == nil {
		return errors.New("trans Is Empty")
	}
	data, _ := json.Marshal(aiis.TransportPayload{
		Method: aiis.TransportMethodResultPatch,
		Data:   status,
	})

	return trans.SendMessage(aiis.SubjectServiceStatus, data)
}

// 设置接收结果回调
func (s *BaseTransport) SetResultHandler(handler ResultHandler) error {
	trans := s.ITransportService
	if trans == nil {
		return errors.New("trans Is Empty")
	}
	subject := SUBJECT_RESULTS
	fn := func(msg *transport.Message) ([]byte, error) {
		var payload aiis.TransportPayload
		data := msg.Body
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}

		body, _ := json.Marshal(payload.Data)
		var result aiis.PublishResult
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		handler(&result)
		return nil, nil
	}

	return trans.OnMessage(subject, fn)
}

func NewBaseTransport(t aiis.ITransportService) *BaseTransport {
	s := &BaseTransport{}
	s.ITransportService = t
	return s
}
