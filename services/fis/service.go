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
)

const (
	LEN_FIS_MO   = 149
	PRS_START	 = 64
	LEN_PR_VALUE = 3
)

type Diagnostic interface {
	Error(msg string, err error)
}

type Service struct {
	Pmon        *pmon.Service
	Odoo        *odoo.Service
	diag        Diagnostic
	configValue atomic.Value
	mtxFile		sync.Mutex
}

func NewService(d Diagnostic, c Config, pmon *pmon.Service) *Service {
	s := &Service{
		diag: d,
		mtxFile: sync.Mutex{},
	}

	s.Pmon = pmon

	s.configValue.Store(c)
	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	s.Pmon.PmonRegistryEvent(s.OnPmonEventMission, s.Config().CHRecvMission, nil)
	s.Pmon.PmonRegistryEvent(s.OnPmonEventHeartbeat, s.Config().CHRecvHeartbeat, nil)
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

	l := len(msg)
	if l != LEN_FIS_MO {
		s.diag.Error(msg, fmt.Errorf("msg len err:%d", l))
		return
	}

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
