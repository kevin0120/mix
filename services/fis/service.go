package fis

import (
	"fmt"
	"github.com/masami10/aiis/services/odoo"
	"github.com/masami10/aiis/services/pmon"
	"strconv"
	"sync/atomic"
	"io/ioutil"
	"strings"
	"github.com/willf/pad"
	"sync"
	"time"
)

const (
	PRS_START	 = 64
	LEN_PR_VALUE = 3
	MAX_HEARTBEAT_CHECK_COUNT = 3
	FIS_STATUS_ONLINE = "online"
	FIS_STATUS_OFFLINE = "offline"
)

type Diagnostic interface {
	Error(msg string, err error)
}

type Service struct {
	Pmon        	*pmon.Service
	Odoo        	*odoo.Service
	diag        	Diagnostic
	configValue 	atomic.Value
	mtxFile			sync.Mutex
	keepAliveCount	int32
	status			string
	mtxStatus		sync.Mutex
}

func NewService(d Diagnostic, c Config, pmon *pmon.Service) *Service {
	s := &Service{
		diag: d,
		mtxFile: sync.Mutex{},
		status: FIS_STATUS_OFFLINE,
		mtxStatus: sync.Mutex{},
	}

	s.Pmon = pmon

	s.configValue.Store(c)
	return s
}

func (s *Service) UpdateStatus(status string) {
	s.mtxStatus.Lock()
	defer s.mtxStatus.Unlock()

	if s.status != status {
		s.status = status
	}
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	c := s.Config()
	s.Pmon.PmonRegistryEvent(s.OnPmonEventMission, c.CHRecvMission, nil)
	s.Pmon.PmonRegistryEvent(s.OnPmonEventHeartbeat, c.CHRecvHeartbeat, nil)

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
		// 处理pmon收到的数据
		msg := string(data)
		s.HandleMO(msg)
	}
}

func (s *Service) OnPmonEventHeartbeat(err error, data []rune, obj interface{}) {
	s.UpdateStatus(FIS_STATUS_ONLINE)
	s.updateKeepAliveCount(0)
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
	// 更新通道restartpoint
	s.Pmon.Channels[ch].RefreshRestartPoint(restartPoint)

	// 保存文件
	s.mtxFile.Lock()
	defer s.mtxFile.Unlock()

	f, err := ioutil.ReadFile(s.Pmon.Config().Ofhkht)
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

			lines[i] = fmt.Sprint("%s*%s*", ch, pad.Left(restartPoint, s.Pmon.Channels[ch].RestartPointLength, "0"))
			break
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(s.Pmon.Config().Ofhkht, []byte(output), 0644)
	if err != nil {
		s.diag.Error("save restart point", err)
	}
}

func (s *Service) HandleMO(msg string) {

	mo := odoo.ODOOMO{}

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
	s.SaveRestartPoint(mo.Lnr, s.Config().CHRecvMission)

	// prs
	num_prs := len(s.Config().PRS)
	prs_end := PRS_START + (LEN_PR_VALUE * num_prs + num_prs - 1)
	s_prs := msg[PRS_START:prs_end]
	var step = 0
	for i := 0; i < num_prs; i++ {
		pr := odoo.ODOOPR{}
		pr.Pr_group = s.Config().PRS[i]
		pr.Pr_value = s_prs[step : step+LEN_PR_VALUE]

		mo.Prs = append(mo.Prs, pr)
		step += LEN_PR_VALUE + 1
	}

	err := s.Odoo.CreateMO(mo)
	if err != nil {
		s.diag.Error("create mo err", err)
	}
}

func (s *Service) PushResult(result *FisResult) error {
	return s.Pmon.SendPmonMessage(pmon.PMONMSGSD, s.Config().CHSendResult, result.Serialize())
}
