package networkModule

import (
	"driverModule"
	"net"
)

const (
	IMA_PERIOD        = 25
	IMA_LOSS          = 4
	ALIVE_WATCH       = 10
	NET_SETUP         = 200
	DIAL_INT          = 50
	CONN_ATMPT        = 5
	WRITE_DL          = 10
	READ_DL           = 10
	CONN_FAIL_TIMEOUT = 2 * NET_SETUP
)

const (
	ORDER_TAKEN           = "OTK"
	ORDER_EXECUTED        = "OEX"
	TAKE_BACKUP_ORDER     = "TBO"
	BACKUP_ORDER_COMPLETE = "BOC"
	TAKE_NEW_ORDER        = "TNO"
	TAKE_BACKUP_FLOOR     = "TBF"
	ENGINE_FAILURE        = "ENF"
	ENGINE_RECOVERY       = "ENR"
)

var (
	BROAD_CAST = GetBroadcastIP()
	LOCAL_IP   = GetLocalIP()
	UDP_PORT   = "9001"
	TCP_PORT   = "9191"
)

type connMap struct {
	tcpMap map[string]connChans
}

type connChans struct {
	send chan Mail
	quit chan bool
}

type tcpConnection struct {
	ip       string
	socket   net.Conn
	sendChan chan Mail
	quit     chan bool
}

type Message struct {
	SendersIP   string
	Type        int
	Destination int
	Floor       int
	Order       int
	GlobalQueue [driverModule.N_FLOORS * 2]int
}

type Mail struct {
	TargetIP string
	Msg      Message
}
