package networkModule

type internalChannels struct {
	setupFail       chan bool
	ima             chan string
	newIP           chan string
	deadElevator    chan string
	updateTcpMap    chan tcpConnection
	connectFail     chan string
	closeConn       chan string
	errorIP         chan string
	connectionError chan string
	quitImaSend     chan bool
	quitImaListen   chan bool
	quitImaWatcher  chan bool
	quitListenTcp   chan bool
	quitTcpMap      chan bool
}

type NetChannels struct {
	GetDeadElevator  chan string
	SendDeadElevator chan string
	SendToAll        chan Mail
	SendToOne        chan Mail
	Inbox            chan Mail
	NumOfPeers       chan int
	NewConnection    chan string
}

var internalChan internalChannels
var externalChan NetChannels

func (internalChan *internalChannels) init() {
	internalChan.setupFail = make(chan bool)
	internalChan.ima = make(chan string)
	internalChan.newIP = make(chan string)
	internalChan.deadElevator = make(chan string)
	internalChan.updateTcpMap = make(chan tcpConnection)
	internalChan.connectFail = make(chan string)
	internalChan.connectionError = make(chan string)
	internalChan.errorIP = make(chan string)
	internalChan.closeConn = make(chan string)
	internalChan.quitImaSend = make(chan bool)
	internalChan.quitImaListen = make(chan bool)
	internalChan.quitImaWatcher = make(chan bool)
	internalChan.quitListenTcp = make(chan bool)
	internalChan.quitTcpMap = make(chan bool)
}

func (externalChan *NetChannels) Init() {
	externalChan.GetDeadElevator = make(chan string)
	externalChan.SendToAll = make(chan Mail)
	externalChan.SendToOne = make(chan Mail)
	externalChan.Inbox = make(chan Mail, 40)
	externalChan.NumOfPeers = make(chan int)
	externalChan.NewConnection = make(chan string, 10)

}
