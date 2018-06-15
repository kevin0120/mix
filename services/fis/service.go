package fis

import (
	"github.com/masami10/aiis/services/pmon"
	"fmt"
	"github.com/masami10/aiis/services/odoo"
	"strconv"
	"sync/atomic"
)

const (
	LEN_FIS_MO = 149
	NUM_PRS = 16
	LEN_PR_VALUE = 3
	MAX_RESULT_SEQ = 999999
)

type Diagnostic interface {
	Error(msg string, err error)
}

type Service struct {
	Pmon		*pmon.Service
	Odoo		*odoo.Service
	PR_GROUPS	[]string
	diag        Diagnostic
	configValue atomic.Value

}

func NewService(d Diagnostic, c Config, pmon *pmon.Service) (*Service) {
	s := &Service{
		diag: d,
	}

	s.Pmon = pmon
	s.PR_GROUPS = []string {
		"GSP",
		"SAB",
		"RAD",
		"REI",
		"BRS",
		"GMO",
		"EDF",
		"AED",
		"MOT",
		"LES",
		"HIS",
		"AIB",
		"FSB",
		"SIH",
		"HBV",
		"AGM",
	}

	s.configValue.Store(c)
	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	s.Pmon.PmonRegistryEvent(s.OnPmonEvent, s.Config().CH_RECV, nil)

	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) OnPmonEvent(err error, data []rune, obj interface{}) {
	if err != nil {
		fmt.Printf("err %s\n", err.Error())
	} else {
		// 处理pmon收到的数据
		msg := string(data)
		s.HandleMO(msg)
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

	// prs
	s_prs := msg[64:127]
	var step = 0
	for i := 0; i < NUM_PRS; i++ {
		pr := odoo.ODOOPR{}
		pr.Pr_group = s.PR_GROUPS[i]
		pr.Pr_value = s_prs[step:step + LEN_PR_VALUE]

		mo.Prs = append(mo.Prs, pr)
		step += LEN_PR_VALUE + 1
	}

	err := s.Odoo.CreateMO(mo)
	if err != nil {
		s.diag.Error("create mo err", err)
	}
}

func (s *Service) PushResult(fis_result *FisResult) (error) {
	return s.Pmon.SendPmonMessage(pmon.PMONMSGSD, s.Config().CH_SEND, fis_result.Serialize())
}