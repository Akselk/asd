package networkModule

import (
	"colorPrint"
	"net"
	"time"
)

func imaListen() {
	service := BROAD_CAST + ":" + UDP_PORT
	addr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		colorPrint.DataWithColor(colorPrint.COLOR_RED, "network.IMAListen()--> ResolveUDP error")
		internalChan.setupFail <- true
	}
	sock, err := net.ListenUDP("udp4", addr)
	if err != nil {
		colorPrint.DataWithColor(colorPrint.COLOR_RED, "network.IMAListen()--> ListenUDP error")
	}
	var data [512]byte
	for {
		select {
		case <-internalChan.quitImaListen:
			return
		default:
			_, remoteAddr, err := sock.ReadFromUDP(data[0:])
			if err != nil {
				colorPrint.DataWithColor(colorPrint.COLOR_RED, "network.IMAListen()--> ReadFromUDP error")
				break
			}
			if LOCAL_IP != remoteAddr.IP.String() {
				if err == nil {
					elevIP := remoteAddr.IP.String()
					internalChan.ima <- elevIP
				} else {
					colorPrint.DataWithColor(colorPrint.COLOR_RED, "network.IMAListen()--> UDP read error")
				}
			}
		}
	}
}

func imaWatcher() {
	peers := make(map[string]time.Time)
	deadline := IMA_LOSS * IMA_PERIOD * time.Millisecond
	for {
		select {
		case ip := <-internalChan.ima:
			_, inMap := peers[ip]
			if inMap {
				peers[ip] = time.Now()
			} else {
				peers[ip] = time.Now()
				internalChan.newIP <- ip
			}
		case <-time.After(ALIVE_WATCH * time.Millisecond):
			for ip, timestamp := range peers {
				if time.Now().After(timestamp.Add(deadline)) {
					colorPrint.DataWithColor(colorPrint.COLOR_RED, "network.imaWatcher --> Timeout", ip)
					externalChan.GetDeadElevator <- ip
					internalChan.closeConn <- ip
					delete(peers, ip)
				}
			}
		case <-internalChan.quitImaWatcher:
			return
		}
	}
}

func imaSend() {
	service := BROAD_CAST + ":" + UDP_PORT
	addr, err := net.ResolveUDPAddr("udp4", service)
	addrSelf, err := net.ResolveUDPAddr("udp", "localhost:20014")
	msg := make([]byte, 1)
	msg[0] = byte(1)
	if err != nil {
		colorPrint.DataWithColor(colorPrint.COLOR_RED, "network.IMASend()--> Resolve error")
		internalChan.setupFail <- true
	}
	imaSock, err := net.DialUDP("udp4", nil, addr)
	imaSockSelf, err := net.DialUDP("udp", nil, addrSelf)
	if err != nil {
		colorPrint.DataWithColor(colorPrint.COLOR_RED, "network.IMASend()--> Dial error")
		internalChan.setupFail <- true
	}
	ima := []byte("IMA")
	for {
		select {
		case <-internalChan.quitImaSend:
			return
		default:
			_, err := imaSock.Write(ima)
			imaSockSelf.Write(msg)

			if err != nil {
				colorPrint.DataWithColor(colorPrint.COLOR_RED, "network.IMASend()--> UDP send error")
			}
			time.Sleep(IMA_PERIOD * time.Millisecond)
		}
	}
}
