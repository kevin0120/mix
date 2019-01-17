package pmon

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	MINPMONBUFSIZE = 256
)

type SegmentType = string

const (
	SEGMENT   = "02"
	NOSEGMENT = "01"
	BOTH      = "03"
)

type ChannelType = int

const (
	PMONMASTER = iota
	PMONSLAVE
)

type Config struct {
	Path           string `yaml:"path"`
	Workers        int    `yaml:"workers"`
	Enable         bool   `yaml:"enable"`
	SendCheckCount int    `yaml:"send_check_count"`
}

type PmonConfig struct {
	ChannelCount  int
	Channels      map[string]cChannel
	WaitResp      time.Duration
	UdpCount      int
	Connections   map[string]cConnection
	Ofhkht        string
	RestartPoints map[string]string
}

type cConnection struct {
	Board   int
	Port    int
	Count   int
	Address []string //一个连接最大可以定义4个地址
}

type cChannel struct {
	Ch                 string
	Area               string
	Base               int           //for PLC
	Offs               int           //for PLC
	Buffer             int           //bytes 数量，最小256， 因为协议头18个字节，所以buffer要定义需要数量+18
	Port               int           //通道赋予那个连接，连接索引
	Type               ChannelType   //M:Master， S:Slave
	SNoT               string        // Number of transmitter
	SNoR               string        // Number of recipient
	Segment            SegmentType   // Y segment, N: nosegment
	WriteTimeout       time.Duration //为了重新发送数据
	RestartPointLength int
	Description        string
}

func parseRestartPoints(conf *PmonConfig) error {
	f, err := os.Open(conf.Ofhkht)
	if err != nil {
		return err
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')

		if line != "" {
			values := strings.Split(line, "*")
			if len(values) < 2 {
				return errors.New("restart point format error")
			}

			conf.RestartPoints[values[0]] = values[1]
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

func parseChannel(key *ini.Key) (string, *cChannel) {
	ch := key.String() // 拿到了一整行作为字符串，解析
	if ch == "" {
		return ch, nil
	}
	sses := strings.Split(ch, ",")

	var ss []string

	for _, s := range sses {
		d := strings.Replace(s, " ", "", -1)
		ss = append(ss, d)
	}

	if len(ss) < 12 {
		return ch, nil
	}

	if ss[1] != "PC" {
		//不是PC端使用，应该抛弃
		fmt.Printf("error config for Area， must be PC, now %s", ss[1])
		return ch, nil
	}
	base, _ := strconv.Atoi(ss[2])
	offset, _ := strconv.Atoi(ss[3])
	buf, _ := strconv.Atoi(ss[4])
	port, _ := strconv.Atoi(ss[5])
	timeout, _ := strconv.Atoi(ss[10])
	rpl, _ := strconv.Atoi(ss[11])
	if buf < MINPMONBUFSIZE {
		//小于协议定义的最小buff大小
		fmt.Printf("error config for buffer size, min: %d, now: %d", MINPMONBUFSIZE, buf)
		return ch, nil
	}
	var seg = BOTH
	if ss[9] == "Y" {
		seg = SEGMENT
	} else if ss[9] == "N" {
		seg = NOSEGMENT
	}

	var sType = PMONMASTER
	if ss[6] == "M" {
		sType = PMONMASTER
	} else if ss[6] == "S" {
		sType = PMONSLAVE
	}
	return ss[0], &cChannel{
		Ch:                 ss[0],
		Area:               ss[1],
		Base:               base,
		Offs:               offset,
		Buffer:             buf,
		Port:               port,
		Type:               sType,
		SNoT:               ss[7],
		SNoR:               ss[8],
		Segment:            seg,
		WriteTimeout:       time.Duration(time.Duration(timeout) * time.Second),
		RestartPointLength: rpl,
		Description:        ss[12],
	}
}

func parseConnection(key *ini.Key) cConnection {
	ch := key.String() // 拿到了一整行作为字符串，解析
	if ch == "" {
		return cConnection{}
	}
	_ss := strings.Split(ch, ",")

	var ss []string

	for _, s := range _ss {
		d := strings.Replace(s, " ", "", -1)
		ss = append(ss, d)
	}

	board, _ := strconv.Atoi(ss[0])
	port, _ := strconv.Atoi(ss[1])
	count, _ := strconv.Atoi(ss[2])
	c := cConnection{
		Board: board,
		Port:  port,
		Count: count,
	}

	for _, add := range ss[3:] {
		c.Address = append(c.Address, add)
	}
	return c
}

func NewConfig() Config {
	return Config{
		Path:           "/etc/pmon/PMON.CFG",
		Workers:        4,
		Enable:         true,
		SendCheckCount: 10,
	}
}

func (c Config) Validate() error {
	return nil
}

func PmonNewConfig(path string) (PmonConfig, error) {
	if !filepath.IsAbs(path) {
		return PmonConfig{}, fmt.Errorf("path : %s is not a absoluted path", path)
	}

	var conf = PmonConfig{ChannelCount: 0, RestartPoints: map[string]string{}}

	cfg, err := ini.Load(path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
	}

	for _, sec := range cfg.Sections() {
		if strings.Contains(sec.Name(), "ChannelCount") {
			cc, _ := sec.Key("channelcount").Int()
			conf.ChannelCount = cc
			conf.Channels = make(map[string]cChannel, cc)
		}
		if strings.Contains(sec.Name(), "WaitResp") {
			cc, _ := sec.Keys()[0].Int()
			conf.WaitResp = time.Duration(time.Millisecond * time.Duration(cc))
		}
		if strings.Contains(sec.Name(), "UdpCount") {
			cc, _ := sec.Keys()[0].Int()
			conf.UdpCount = cc
		}
		if strings.Contains(sec.Name(), "Channels") {
			for _, key := range sec.Keys() {
				c, ch := parseChannel(key)
				if ch != nil {
					conf.Channels[c] = *ch
				}
			}
		}
		if strings.Contains(sec.Name(), "Connections") {
			conf.Connections = make(map[string]cConnection, len(sec.Keys()))
			for _, key := range sec.Keys() {
				con := parseConnection(key)
				conf.Connections[key.Name()] = con
			}
		}
		if strings.Contains(sec.Name(), "RestartPoint") {
			for _, key := range sec.Keys() {
				conf.Ofhkht = key.Value()
				break
			}

			err := parseRestartPoints(&conf)
			if err != nil {
				fmt.Printf("parseRestartPoints failed: %v\n", err)
			}
		}
	}

	return conf, nil
}
