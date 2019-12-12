package wsnotify

import "github.com/kataras/iris/websocket"

type DispatcherNotifyPackage struct {
	C    websocket.Connection
	Data []byte
}
