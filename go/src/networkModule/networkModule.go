package networkModule

import (
	"colorPrint"
	"driverModule"
	"math/rand"
	"time"
)

func NetworkSetup(NetChan NetChannels) {
	rand.Seed(time.Now().UTC().UnixNano())

	internalChan.init()

	externalChan = NetChan

	imaStart()
	networkStart()
}

func imaStart() {
	go imaWatcher()
	go imaListen()
	go imaSend()
}

func networkStart() {
	go manageTcpConnections()

	for {
		select {

		case <-internalChan.setupFail:
			colorPrint.DataWithColor(colorPrint.COLOR_RED, "net.Startup--> Setupfail. Retrying...")

			internalChan.quitImaSend <- true
			internalChan.quitImaListen <- true
			internalChan.quitImaWatcher <- true
			internalChan.quitListenTcp <- true
			internalChan.quitTcpMap <- true
			time.Sleep(time.Millisecond)

			imaStart()

			go manageTcpConnections()

		case <-time.After(NET_SETUP * time.Millisecond):
			return
		}
	}
}

func (mail *Mail) MakeMail(TargetIP string, MyIP string, Type int, currDest int, floor int, order int, queue [2 * driverModule.N_FLOORS]int) {

	mail.TargetIP = TargetIP
	mail.Msg.SendersIP = MyIP
	mail.Msg.Type = Type
	mail.Msg.Destination = currDest
	mail.Msg.Floor = floor
	mail.Msg.Order = order
	mail.Msg.GlobalQueue = queue
}
