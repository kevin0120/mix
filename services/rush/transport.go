package rush

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/transport"
)

type ResultHandler func(clientID string, result *aiis.PublishResult)

type BaseTransport struct {
	aiis.BaseTransport
}

// 发送结果上传反馈
func (s *BaseTransport) SendResultPatch(clientID string, rp *aiis.ResultPatch) error {
	trans := s.ITransportService
	if trans == nil {
		return errors.New("trans Is Empty")
	}
	data, _ := json.Marshal(aiis.TransportPayload{
		Method: aiis.TransportMethodResultPatch,
		Data:   rp,
	})

	return trans.SendMessage(fmt.Sprintf(SUBJECT_RESULTS_RESP, clientID), data)
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

func getMsgHeaderProperty(message *transport.Message, property string) string {
	header := message.Header
	pp := property
	switch property {
	case transport.HEADER_SUBJECT, transport.HEADER_REPLY, transport.HEADER_MSG_ID, transport.HEADER_CLIENT_ID:
		pp = property
	default:
		return ""
	}
	if ret, ok := header[pp]; ok {
		return ret
	} else {
		return ""
	}
}

func getClientID(message *transport.Message) string {
	return getMsgHeaderProperty(message, transport.HEADER_CLIENT_ID)
}

// 设置接收结果回调
func (s *BaseTransport) SetResultHandler(handler ResultHandler) error {
	trans := s.ITransportService
	if trans == nil {
		return errors.New("trans Is Empty")
	}
	subject := SUBJECT_RESULTS
	fn := func(msg *transport.Message) ([]byte, error) {
		clientID := getClientID(msg)
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

		handler(clientID, &result)
		return nil, nil
	}

	return trans.OnMessage(subject, fn)
}

func NewBaseTransport(t aiis.ITransportService) *BaseTransport {
	s := &BaseTransport{}
	s.ITransportService = t
	return s
}
