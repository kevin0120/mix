package openprotocol

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_parseOpenProtocolErrorCode(t *testing.T) {
	err := "111"
	ss := parseOpenProtocolErrorCode(err)
	t.Log(ss)
}

func procHandleRecv() {
	//c.receiveBuf = make(chan []byte, BufferSize)
	//handleRecvBuf := make([]byte, BufferSize)
	//lenBuf := len(handleRecvBuf)
	handleRecvBuf := []byte{}
	var writeOffset = 0
	buf := []byte{1, 2, 3, 4, 5, 6, 7, OpTerminal, 1, 2, 3, 4, 5, 6, 7, OpTerminal}

	for {
		// 处理接收缓冲
		var readOffset = 0
		var index = 0

		for {
			if readOffset >= len(buf) {
				break
			}
			index = bytes.IndexByte(buf[readOffset:], OpTerminal)
			if index == -1 {
				// 没有结束字符,放入缓冲等待后续处理
				restBuf := buf[readOffset:]
				//if writeOffset+len(restBuf) > lenBuf {
				//	c.diag.Error("full", errors.New("full"))
				//	break
				//}

				copy(handleRecvBuf[writeOffset:writeOffset+len(restBuf)], restBuf)
				writeOffset += len(restBuf)
				break
			} else {
				// 找到结束字符，结合缓冲进行处理
				targetBuf := append(handleRecvBuf[0:writeOffset], buf[readOffset:readOffset+index]...)
				fmt.Println(targetBuf)
				//err := c.handlePackageOPPayload(targetBuf)
				//if err != nil {
				//	//数据需要丢弃
				//	//c.diag.Error("handlePackageOPPayload Error", err)
				//	//c.diag.Debug(fmt.Sprintf("procHandleRecv Raw Msg:%s", string(buf)))
				//	//c.diag.Debug(fmt.Sprintf("procHandleRecv Rest Msg:%s writeOffset:%d readOffset:%d index:%d", string(handleRecvBuf), writeOffset, readOffset, index))
				//}

				writeOffset = 0
				readOffset += index + 1
			}
		}
	}
}

func Test_HandleBuf(t *testing.T) {
	procHandleRecv()
}
