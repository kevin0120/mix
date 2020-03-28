package ts002

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/services/wsnotify"
)

func (s *Service) OnWS_TOOL_ENABLE(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	req := tightening_device.ToolControl{}
	_ = json.Unmarshal(byteData, &req)

	locationBody, err := s.storage.GetToolLocation(req.ToolSN)
	if err != nil {
		s.diag.Error(fmt.Sprintf("GetToolLocation Failed: %s", req.ToolSN), err)
		return
	}

	location := Location{}
	_ = json.Unmarshal([]byte(locationBody), &location)
	var status uint16 = 0
	if req.Enable {
		status = 1
	}

	if err := s.ioDoAction(location.EquipmentSN, location.Output, status); err != nil {
		s.diag.Error(fmt.Sprintf("ioDoAction Failed: %s", req.ToolSN), err)
		return
	}
}
