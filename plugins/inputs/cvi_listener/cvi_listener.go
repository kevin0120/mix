package cvi_listener

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/masami10/rush"
	"github.com/masami10/rush/internal"
	"github.com/masami10/rush/plugins/inputs"
	"github.com/masami10/rush/plugins/parsers"
	"github.com/masami10/rush/utils/cvi"
	"encoding/xml"
)

type setReadBufferer interface {
	SetReadBuffer(bytes int) error
}

type streamSocketListener struct {
	net.Listener
	*CVIListener

	connections    map[string]net.Conn
	connectionsMtx sync.Mutex

	cvi3_manager 	*CVI3Manager
	api_server		*ApiServer
}

func (ssl *streamSocketListener) listen() {

	ssl.connections = map[string]net.Conn{}

	for {
		c, err := ssl.Accept()
		if err != nil {
			if !strings.HasSuffix(err.Error(), ": use of closed network connection") {
				ssl.AddError(err)
			}
			break
		}

		ssl.connectionsMtx.Lock()
		if ssl.MaxConnections > 0 && len(ssl.connections) >= ssl.MaxConnections {
			ssl.connectionsMtx.Unlock()
			c.Close()
			continue
		}
		ssl.connections[c.RemoteAddr().String()] = c
		ssl.connectionsMtx.Unlock()

		if err := ssl.setKeepAlive(c); err != nil {
			ssl.AddError(fmt.Errorf("unable to configure keep alive (%s): %s", ssl.ServiceAddress, err))
		}

		go ssl.read(c)
	}

	ssl.connectionsMtx.Lock()
	for _, c := range ssl.connections {
		c.Close()
	}
	ssl.connectionsMtx.Unlock()
}

func (ssl *streamSocketListener) setKeepAlive(c net.Conn) error {
	if ssl.KeepAlivePeriod == nil {
		return nil
	}
	tcpc, ok := c.(*net.TCPConn)
	if !ok {
		return fmt.Errorf("cannot set keep alive on a %s socket", strings.SplitN(ssl.ServiceAddress, "://", 2)[0])
	}
	if ssl.KeepAlivePeriod.Duration == 0 {
		return tcpc.SetKeepAlive(false)
	}
	if err := tcpc.SetKeepAlive(true); err != nil {
		return err
	}
	return tcpc.SetKeepAlivePeriod(ssl.KeepAlivePeriod.Duration)
}

func (ssl *streamSocketListener) removeConnection(c net.Conn) {
	defer ssl.connectionsMtx.Unlock()

	ssl.connectionsMtx.Lock()
	delete(ssl.connections, c.RemoteAddr().String())

}

func (ssl *streamSocketListener) read(c net.Conn) {
	defer ssl.removeConnection(c)
	defer c.Close()

	//scnr := bufio.NewScanner(c)
	buffer := make([]byte, 65535)
	for {
		if ssl.ReadTimeout != nil && ssl.ReadTimeout.Duration > 0 {
			c.SetReadDeadline(time.Now().Add(ssl.ReadTimeout.Duration))
		}
		//if !scnr.Scan() {
		//	break
		//}
		//
		//msg := string(scnr.Bytes())
		//fmt.Printf("%s\n", msg)
		n, err := c.Read(buffer)
		if err != nil {
			break
		}
		msg := string(buffer[0:n])
		//fmt.Printf("%s\n", msg)

		header := cvi.CVI3Header{}
		header.Deserialize(msg[0: cvi.HEADER_LEN])
		var body string = msg[cvi.HEADER_LEN: n]
		var rest int = int(header.SIZ) - cvi.HEADER_LEN - n
		for {
			if rest <= 0 {
				break
			}
			n, err := c.Read(buffer)
			if err != nil {
				break
			}
			body += string(buffer[0:n])
			rest -= n
		}

		//fmt.Printf("%s\n", body)
		go ssl.handle(body)

		if header.TYP == cvi.Header_type_request_with_reply || header.TYP == cvi.Header_type_keep_alive {
			// 执行应答
			reply := cvi.CVI3Header{}
			reply.Init()
			reply.TYP = cvi.Header_type_reply
			reply.MID = header.MID
			reply_packet := reply.Serialize()

			_, err := c.Write([]byte(reply_packet))
			if err != nil {
				print("%s\n", err.Error())
				break
			}
		}
	}

	//if err := scnr.Err(); err != nil {
	//	if err, ok := err.(net.Error); ok && err.Timeout() {
	//		log.Printf("D! Timeout in plugin [input.socket_listener]: %s", err)
	//	} else if !strings.HasSuffix(err.Error(), ": use of closed network connection") {
	//		ssl.AddError(err)
	//	}
	//}
}

func (ssl *streamSocketListener) handle(body string) {
	//test_body, e := ioutil.ReadFile("/home/linshenqi/Documents/cur_test.xml")
	//if e != nil {
	//	fmt.Printf("%s\n", e.Error())
	//}
	//result := cvi.CVI3Result{}
	////reader := strings.NewReader(string(test_body))
	//err := xml.Unmarshal(test_body, &result)
	//if err != nil {
	//	fmt.Printf("%s\n", err.Error())
	//}
	//
	//bs, _ :=json.Marshal(result)
	//fmt.Printf("%s\n", string(bs))
	//
	//i := 3
	//i++

	if strings.Contains(body, cvi.XML_RESULT_KEY) {
		//fmt.Printf("%s\n", body)

		result := cvi.CVI3Result{}
		err := xml.Unmarshal([]byte(body), &result)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		//i := 3
		//i++
	}

}

type packetSocketListener struct {
	net.PacketConn
	*CVIListener
}

func (psl *packetSocketListener) listen() {
	buf := make([]byte, 64*1024) // 64kb - maximum size of IP packet
	for {
		n, _, err := psl.ReadFrom(buf)
		if err != nil {
			if !strings.HasSuffix(err.Error(), ": use of closed network connection") {
				psl.AddError(err)
			}
			break
		}

		metrics, err := psl.Parse(buf[:n])
		if err != nil {
			psl.AddError(fmt.Errorf("unable to parse incoming packet: %s", err))
			//TODO rate limit
			continue
		}
		for _, m := range metrics {
			psl.AddFields(m.Name(), m.Fields(), m.Tags(), m.Time())
		}
	}
}

type (
	CVIListener struct {
		ServiceAddress  string
		MaxConnections  int
		ReadBufferSize  int
		ReadTimeout     *internal.Duration
		KeepAlivePeriod *internal.Duration
		Controllers	[]*CVIConfig

		parsers.Parser
		rush.Accumulator
		io.Closer
	}

)

func (sl *CVIListener) Description() string {
	return "Generic socket listener capable of handling multiple socket types."
}

func (sl *CVIListener) SampleConfig() string {
	return `
  ## URL to listen on
  # service_address = "tcp://:4710"

  ## Maximum number of concurrent connections.
  ## Only applies to stream sockets (e.g. TCP).
  ## 0 (default) is unlimited.
  # max_connections = 1024

  ## Read timeout.
  ## Only applies to stream sockets (e.g. TCP).
  ## 0 (default) is unlimited.
  # read_timeout = "30s"

  ## Maximum socket buffer size in bytes.
  ## For stream sockets, once the buffer fills up, the sender will start backing up.
  ## For datagram sockets, once the buffer fills up, metrics will start dropping.
  ## Defaults to the OS default.
  # read_buffer_size = 65535

  ## Period between keep alive probes.
  ## Only applies to TCP sockets.
  ## 0 disables keep alive probes.
  ## Defaults to the OS configuration.
  # keep_alive_period = "5m"

  ## Data format to consume.
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/masami10/rush/blob/master/docs/DATA_FORMATS_INPUT.md
  # data_format = "influx"

  ## CVI3 Controllers, 
  # [[inputs.cvi_listener.controllers]]
  # 	sn = "1"
  # 	ip = "192.168.1.200"
  # 	port = 4700
  #		hmi = "http://127.0.0.1:8000"
`
}

func (sl *CVIListener) Gather(_ rush.Accumulator) error {
	return nil
}

func (sl *CVIListener) SetParser(parser parsers.Parser) {
	sl.Parser = parser
}

func (sl *CVIListener) Start(acc rush.Accumulator) error {
	sl.Accumulator = acc
	spl := strings.SplitN(sl.ServiceAddress, "://", 2)
	if len(spl) != 2 {
		return fmt.Errorf("invalid service address: %s", sl.ServiceAddress)
	}

	if spl[0] == "unix" || spl[0] == "unixpacket" || spl[0] == "unixgram" {
		// no good way of testing for "file does not exist".
		// Instead just ignore error and blow up when we try to listen, which will
		// indicate "address already in use" if file existed and we couldn't remove.
		os.Remove(spl[1])
	}

	switch spl[0] {
	case "tcp", "tcp4", "tcp6", "unix", "unixpacket":
		l, err := net.Listen(spl[0], spl[1])
		if err != nil {
			return err
		}

		if sl.ReadBufferSize > 0 {
			if srb, ok := l.(setReadBufferer); ok {
				srb.SetReadBuffer(sl.ReadBufferSize)
			} else {
				log.Printf("W! Unable to set read buffer on a %s socket", spl[0])
			}
		}

		ssl := &streamSocketListener{
			Listener:    l,
			CVIListener: sl,
		}

		sl.Closer = ssl
		go ssl.listen()

		ssl.cvi3_manager = &CVI3Manager{}
		ssl.api_server = &ApiServer{}

		// 启动客户端服务
		ssl.cvi3_manager.StartService(sl.Controllers)

		// 启动api服务
		go ssl.api_server.StartService(ssl.cvi3_manager)

		//ssl.handle("")

	case "udp", "udp4", "udp6", "ip", "ip4", "ip6", "unixgram":
		pc, err := net.ListenPacket(spl[0], spl[1])
		if err != nil {
			return err
		}

		if sl.ReadBufferSize > 0 {
			if srb, ok := pc.(setReadBufferer); ok {
				srb.SetReadBuffer(sl.ReadBufferSize)
			} else {
				log.Printf("W! Unable to set read buffer on a %s socket", spl[0])
			}
		}

		psl := &packetSocketListener{
			PacketConn:  pc,
			CVIListener: sl,
		}

		sl.Closer = psl
		go psl.listen()
	default:
		return fmt.Errorf("unknown protocol '%s' in '%s'", spl[0], sl.ServiceAddress)
	}

	if spl[0] == "unix" || spl[0] == "unixpacket" || spl[0] == "unixgram" {
		sl.Closer = unixCloser{path: spl[1], closer: sl.Closer}
	}

	return nil
}

func (sl *CVIListener) Stop() {
	if sl.Closer != nil {
		sl.Close()
		sl.Closer = nil
	}
}

func newSocketListener() *CVIListener {
	parser, _ := parsers.NewInfluxParser()

	return &CVIListener{
		Parser: parser,
	}
}

type unixCloser struct {
	path   string
	closer io.Closer
}

func (uc unixCloser) Close() error {
	err := uc.closer.Close()
	os.Remove(uc.path) // ignore error
	return err
}

func init() {
	inputs.Add("cvi_listener", func() rush.Input { return newSocketListener() })
}
