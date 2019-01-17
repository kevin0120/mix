package fis

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/aiis/services/odoo"
	"github.com/masami10/aiis/services/pmon"
	"github.com/willf/pad"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	PRS_START                 = 64
	LEN_PR_VALUE              = 3
	LEN_MO_TAIL               = 22
	MAX_HEARTBEAT_CHECK_COUNT = 3
	FIS_STATUS_ONLINE         = "online"
	FIS_STATUS_OFFLINE        = "offline"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	Pmon           *pmon.Service
	Odoo           *odoo.Service
	diag           Diagnostic
	configValue    atomic.Value
	mtxFile        sync.Mutex
	keepAliveCount int32
	status         string
	mtxStatus      sync.Mutex
	missionBuf     chan string
	urgRequest     map[string]string
	mtxUrg         sync.Mutex
	resultsBuf     chan string
}

func NewService(d Diagnostic, c Config, pmon *pmon.Service) *Service {
	if c.Enable {
		s := &Service{
			diag:       d,
			mtxFile:    sync.Mutex{},
			status:     FIS_STATUS_OFFLINE,
			mtxStatus:  sync.Mutex{},
			missionBuf: make(chan string, 2048),
			urgRequest: map[string]string{},
			mtxUrg:     sync.Mutex{},
			resultsBuf: make(chan string, 2048),
		}

		s.Pmon = pmon

		s.configValue.Store(c)
		return s
	}

	return nil
}

func (s *Service) manageResults() {
	for {
		select {
		case data := <-s.resultsBuf:
			s.Pmon.SendPmonMessage(pmon.PMONMSGSD, s.Config().CHSendResult, data)
			time.Sleep(time.Duration(s.Config().ADTimeout))
		}
	}
}

func (s *Service) UpdateUrg(longpin string, result string) {
	defer s.mtxUrg.Unlock()
	s.mtxUrg.Lock()

	_, e := s.urgRequest[longpin]
	if e {
		return
	}

	s.urgRequest[longpin] = result
}

func (s *Service) GetUrg(longpin string) string {
	defer s.mtxUrg.Unlock()
	s.mtxUrg.Lock()

	result, e := s.urgRequest[longpin]
	if e {
		return result
	}

	return ""
}

func (s *Service) RemoveUrg(longpin string) {
	defer s.mtxUrg.Unlock()
	s.mtxUrg.Lock()

	_, e := s.urgRequest[longpin]
	if e {
		delete(s.urgRequest, longpin)
	}
}

func (s *Service) UpdateStatus(status string) {
	s.mtxStatus.Lock()
	defer s.mtxStatus.Unlock()

	if s.status != status {
		s.status = status
		s.diag.Debug(fmt.Sprintf("fis status:%s\n", status))
	}
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {

	c := s.Config()
	if len(c.CHRecvHeartbeat) > 0 {
		e := s.Pmon.SendPmonMessage(pmon.PMONMSGSO, c.CHRecvHeartbeat, "")
		if e != nil {
			s.diag.Error("Send PMON SO msg fail", e)
			return e
		}
	}
	if len(c.CHRecvMission) > 0 {
		e := s.Pmon.SendPmonMessage(pmon.PMONMSGSO, c.CHRecvMission, "")
		if e != nil {
			s.diag.Error("Send PMON SO msg fail", e)
			return e
		}
	}
	//if len(c.CHSendResult) > 0 {
	//	e := s.Pmon.SendPmonMessage(pmon.PMONMSGSO, c.CHSendResult, "")
	//	if e != nil {
	//		s.diag.Error("Send PMON SO msg fail", e)
	//		return e
	//	}
	//}
	//
	//if len(c.CHSendUrg) > 0 {
	//	e := s.Pmon.SendPmonMessage(pmon.PMONMSGSO, c.CHSendUrg, "")
	//	if e != nil {
	//		s.diag.Error("Send PMON SO msg fail", e)
	//		return e
	//	}
	//}

	if len(c.CHRecvUrg) > 0 {
		e := s.Pmon.SendPmonMessage(pmon.PMONMSGSO, c.CHRecvUrg, "")
		if e != nil {
			s.diag.Error("Send PMON SO msg fail", e)
			return e
		}
	}

	s.Pmon.PmonRegistryEvent(s.OnPmonEventMission, c.CHRecvMission, nil)
	s.Pmon.PmonRegistryEvent(s.OnPmonEventHeartbeat, c.CHRecvHeartbeat, nil)
	s.Pmon.PmonRegistryEvent(s.OnPmonEventUrg, c.CHRecvUrg, nil)

	go s.TaskMission()
	go s.manageResults()
	go s.HeartbeatCheck()

	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) OnPmonEventMission(err error, data []rune, obj interface{}) {

	if err != nil {
		s.diag.Error("pmon event mission error", err)
	} else {

		s.missionBuf <- string(data)
	}
}

func (s *Service) OnPmonEventHeartbeat(err error, data []rune, obj interface{}) {
	s.UpdateStatus(FIS_STATUS_ONLINE)
	s.updateKeepAliveCount(0)
}

func (s *Service) OnPmonEventUrg(err error, data []rune, obj interface{}) {
	if err != nil {
		s.diag.Error("pmon event urg error", err)
	} else {

		msg := string(data)
		s.diag.Debug(fmt.Sprintf("receive urg mission: %s\n", msg))
		longpin, handleErr := s.HandleMO(msg)
		if handleErr != nil {
			s.UpdateUrg(longpin, FIS_URG_FAIL)
			return
		}

		s.UpdateUrg(longpin, FIS_URG_SUCCESS)
	}
}

func (s *Service) KeepAliveCount() int32 {
	return atomic.LoadInt32(&s.keepAliveCount)
}

func (s *Service) updateKeepAliveCount(i int32) {
	atomic.SwapInt32(&s.keepAliveCount, i)
}

func (s *Service) addKeepAliveCount() {
	atomic.AddInt32(&s.keepAliveCount, 1)
}

func (s *Service) HeartbeatCheck() {
	interval := s.Config().HeartbeatItv
	for {
		select {
		case <-time.After(time.Duration(interval)):
			if s.KeepAliveCount() >= MAX_HEARTBEAT_CHECK_COUNT {
				s.UpdateStatus(FIS_STATUS_OFFLINE)
			}

			s.addKeepAliveCount()

		}
	}
}

func (s *Service) SaveRestartPoint(restartPoint string, ch string) {

	c := s.Pmon.Config()
	// 更新通道restartpoint
	s.Pmon.Channels[ch].RefreshRestartPoint(restartPoint)

	// 保存文件
	s.mtxFile.Lock()
	defer s.mtxFile.Unlock()

	f, err := ioutil.ReadFile(c.Ofhkht)
	if err != nil {
		s.diag.Error("read restart point file err", err)
	}

	lines := strings.Split(string(f), "\n")

	for i, line := range lines {
		values := strings.Split(line, "*")
		if len(values) < 2 {
			s.diag.Error("restart point format error", nil)
		}

		if values[0] == ch {

			lines[i] = fmt.Sprintf("%s*%s*", ch, pad.Left(restartPoint, s.Pmon.Channels[ch].RestartPointLength, "0"))
			break
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(c.Ofhkht, []byte(output), 0644)
	if err != nil {
		s.diag.Error("save restart point", err)
	}
}

func (s *Service) TaskMission() {
	for {
		select {
		case msg := <-s.missionBuf:
			s.diag.Debug(fmt.Sprintf("receive fis mission: %s\n", msg))
			s.HandleMO(msg)

			time.Sleep(time.Duration(s.Config().MissionItv))
		}
	}
}

func (s *Service) HandleMO(str string) (string, error) {

	msg := strings.TrimSpace(str)

	c := s.Config()

	numPrs := len(c.PRS)
	prsEnd := PRS_START + (LEN_PR_VALUE*numPrs + numPrs - 1)
	if len(msg) < prsEnd {
		return "", errors.New("prs len error")
	}
	sPrs := msg[PRS_START:prsEnd]

	len_mo := prsEnd + LEN_MO_TAIL
	if len_mo != len(msg) {
		s.diag.Error(msg, fmt.Errorf("msg len err:%d", len(msg)))
		return "", errors.New("msg len error")
	}

	var mo odoo.ODOOMO

	// 设备名
	mo.Equipment_name = msg[0:4]

	// 工厂代号
	mo.Factory_name = msg[12:14]

	// 年份
	mo.Year, _ = strconv.Atoi(msg[15:19])

	// 装配代码
	mo.Pin, _ = strconv.Atoi(msg[20:27])

	// 装配代码校验位
	mo.Pin_check_code, _ = strconv.Atoi(msg[28:29])

	// 流水线
	mo.Assembly_line = msg[30:32]

	// 车型
	mo.Model = msg[58:64]

	// 车辆识别号
	mo.Vin = msg[40:57]

	// 流水号
	mo.Lnr = msg[34:38]
	s.SaveRestartPoint(mo.Lnr, c.CHRecvMission)

	// prs
	var step = 0
	for i := 0; i < numPrs; i++ {
		pr := odoo.ODOOPR{}
		pr.Pr_group = s.Config().PRS[i]
		pr.Pr_value = sPrs[step : step+LEN_PR_VALUE]

		if strings.TrimSpace(pr.Pr_value) != "" {
			mo.Prs = append(mo.Prs, pr)
		}

		step += LEN_PR_VALUE + 1
	}

	longpin := fmt.Sprint("%s%d%d%d", mo.Factory_name, mo.Year, mo.Pin, mo.Pin_check_code)

	err := s.Odoo.CreateMO(mo)
	if err != nil {
		if strings.Contains(err.Error(), "already") {
			return longpin, nil
		}

		return "", err
		s.diag.Error(str, err)
	}

	return longpin, nil
}

func (s *Service) PushResult(result *FisResult) error {
	s.resultsBuf <- result.Serialize()
	return nil
}

func (s *Service) PushUrgRequest(urg *FisUrgRequest) error {
	urg.EquipemntName = s.Config().EquipmentName
	return s.Pmon.SendPmonMessage(pmon.PMONMSGSD, s.Config().CHSendUrg, urg.Serialize())
}
