package networkModule

import (
	"colorPrint"
	"encoding/json"
	"math/rand"
	"net"
	"strings"
	"time"
)

func manageTcpConnections() {
	connections := connMap{make(map[string]connChans)}
	go listenForTcpConnections()
	for {
		select {
		case newIP := <-internalChan.newIP:
			connections.handleNewIP(newIP)

		case newTcpConnection := <-internalChan.updateTcpMap:
			connections.handleNewConnection(newTcpConnection)

		case errorIP := <-internalChan.connectFail:
			connections.handleFailedToConnect(errorIP)

		case errorIP := <-internalChan.connectionError:
			connections.handleConnectionError(errorIP)

		case closeIP := <-internalChan.closeConn:
			connections.handleCloseConnection(closeIP)

		case mail := <-externalChan.SendToAll:
			connections.handleSendToAll(mail)

		case mail := <-externalChan.SendToOne:
			connections.handleSendToOne(mail)

		case <-internalChan.quitTcpMap:
			return
		}
	}
}

func (connections *connMap) handleNewIP(newIP string) {
	_, inMap := connections.tcpMap[newIP]
	if !inMap {
		go connectTCP(newIP)
	} else {
		colorPrint.DataWithColor(colorPrint.COLOR_YELLOW, "network.monitorTCPConnections-->", newIP, "already in connections")
	}
}

func (connections *connMap) handleNewConnection(conn tcpConnection) {
	_, inMap := connections.tcpMap[conn.ip]
	if !inMap {
		connections.tcpMap[conn.ip] = connChans{send: make(chan Mail), quit: make(chan bool)}
		colorPrint.DataWithColor(colorPrint.COLOR_GREEN, "network.monitorTCPConnections---> Connection made to ", conn.ip)
		conn.sendChan = connections.tcpMap[conn.ip].send
		conn.quit = connections.tcpMap[conn.ip].quit
		externalChan.NewConnection <- conn.ip
		go conn.handleConnection()
		go peerUpdate(len(connections.tcpMap))
	} else {
		colorPrint.DataWithColor(colorPrint.COLOR_YELLOW, "network.monitorTCPConnections--> A connection already exist to", conn.ip)
		conn.socket.Close()
	}
}

func (connections *connMap) handleFailedToConnect(errorIP string) {
	_, inMap := connections.tcpMap[errorIP]
	if inMap {
		colorPrint.DataWithColor(colorPrint.COLOR_YELLOW, "network.monitorTCPConnections--> Could not dial up ", errorIP, "but a connection already exist")
	} else {
		colorPrint.DataWithColor(colorPrint.COLOR_RED, "network.monitorTCPConnections--> Could not connect to ", errorIP)
		internalChan.errorIP <- errorIP
	}
}

func (connections *connMap) handleConnectionError(errorIP string) {
	_, inMap := connections.tcpMap[errorIP]
	if inMap {
		delete(connections.tcpMap, errorIP)
	}
	go connectTCP(errorIP)
}

func (connections *connMap) handleCloseConnection(closeIP string) {
	connChans, inMap := connections.tcpMap[closeIP]
	if inMap {
		select {
		case connChans.quit <- true:
		case <-time.After(10 * time.Millisecond):
		}
		delete(connections.tcpMap, closeIP)
		numOfConns := len(connections.tcpMap)
		if numOfConns == 0 {
			go peerUpdate(numOfConns)
		}
	}
}

func (connections *connMap) handleSendToOne(mail Mail) {
	switch mail.TargetIP {
	case "":
		size := len(connections.tcpMap)
		if size != 0 {
			for _, connChans := range connections.tcpMap {
				connChans.send <- mail
				break
			}
		}
	default:
		connChans, inMap := connections.tcpMap[mail.TargetIP]
		if inMap {
			connChans.send <- mail
		} else {
			internalChan.errorIP <- mail.TargetIP
		}
	}
}

func (connections *connMap) handleSendToAll(mail Mail) {
	if len(connections.tcpMap) != 0 {
		for _, connChans := range connections.tcpMap {
			connChans.send <- mail
		}
	}
}

func (conn *tcpConnection) handleConnection() {
	quitInbox := make(chan bool)
	go conn.inbox(quitInbox)
	connectionEncoder := json.NewEncoder(conn.socket)
	for {
		select {
		case mail := <-conn.sendChan:
			encodedMsg := mail.Msg
			err := connectionEncoder.Encode(&encodedMsg)
			if err == nil {
			} else {
				colorPrint.DataWithColor(colorPrint.COLOR_RED, "Network.handleConnection--> Error sending message to ", conn.ip, err)
				internalChan.connectionError <- conn.ip
			}
		case <-conn.quit:
			conn.socket.Close()
			colorPrint.DataWithColor(colorPrint.COLOR_YELLOW, "Network.handleConnections--> Connection to ", conn.ip, " has been terminated.")
			return
		case <-quitInbox:
			conn.socket.Close()
			colorPrint.DataWithColor(colorPrint.COLOR_YELLOW, "Network.handleConnections--> Connection to ", conn.ip, " has been terminated.")
			internalChan.connectionError <- conn.ip
			return
		}
	}
}

func (conn *tcpConnection) inbox(quitInbox chan bool) {
	connectionDecoder := json.NewDecoder(conn.socket)
	for {
		decodedMsg := new(Message)
		err := connectionDecoder.Decode(decodedMsg)
		switch err {
		case nil:
			newMail := Mail{TargetIP: conn.ip, Msg: *decodedMsg}
			externalChan.Inbox <- newMail
		default:
			colorPrint.DataWithColor(colorPrint.COLOR_RED, "Network.inbox--> Error:", err)
			time.Sleep(IMA_PERIOD * IMA_LOSS * 2 * time.Millisecond)
			select {
			case quitInbox <- true:
			case <-time.After(WRITE_DL * time.Millisecond):
			}
			return
		}
	}
}

func connectTCP(ip string) {
	attempts := 0
	for attempts < CONN_ATMPT {
		service := ip + ":" + TCP_PORT
		_, err := net.ResolveTCPAddr("tcp4", service)
		if err != nil {
			colorPrint.DataWithColor(colorPrint.COLOR_RED, "Network.connectTCP--> ResolveTCPAddr failed")
			attempts++
			time.Sleep(DIAL_INT * time.Millisecond)
		} else {
			randSleep := time.Duration(rand.Intn(500)+500) * time.Microsecond
			time.Sleep(randSleep)
			socket, err := net.Dial("tcp4", service)
			if err != nil {
				colorPrint.DataWithColor(colorPrint.COLOR_RED, "Network.connectTCP--> DialTCP error when connecting to", ip, " error: ", err)
				attempts++
				time.Sleep(DIAL_INT * time.Millisecond)
			} else {
				newTcpConnection := tcpConnection{ip: ip, socket: socket}
				internalChan.updateTcpMap <- newTcpConnection
				break
			}
		}
	}
	if attempts == CONN_ATMPT {
		select {
		case internalChan.connectFail <- ip:
		case <-time.After(CONN_FAIL_TIMEOUT * time.Millisecond):
			return
		}
	}
}

func listenForTcpConnections() {
	service := ":" + TCP_PORT
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		colorPrint.DataWithColor(colorPrint.COLOR_RED, "Network.listenForTcpConnections--> TCP resolve error")
		internalChan.setupFail <- true
	} else {
		listenSock, err := net.ListenTCP("tcp4", tcpAddr)
		if err != nil {
			colorPrint.DataWithColor(colorPrint.COLOR_RED, "Network.connectTCP--> ListenTCP error")
			internalChan.setupFail <- true
		} else {
			for {
				select {
				case <-internalChan.quitListenTcp:
					return
				default:
					socket, err := listenSock.Accept()
					if err == nil {
						ip := cleanUpIP(socket.RemoteAddr().String())
						newTcpConnection := tcpConnection{ip: ip, socket: socket}
						internalChan.updateTcpMap <- newTcpConnection
					}
				}
			}
		}
	}
}

func peerUpdate(NumOfPeers int) {
	select {
	case externalChan.NumOfPeers <- NumOfPeers:
	case <-time.After(500 * time.Millisecond):
	}
}

func cleanUpIP(garbage string) (cleanIP string) {
	split := strings.Split(garbage, ":")
	cleanIP = split[0]
	return
}
