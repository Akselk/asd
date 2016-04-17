package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"
)

func backup(udpListen *net.UDPConn) {
	listenChan := make(chan int, 1)
	crashcounter := 0
	go listen(listenChan, udpListen)
	for {
		select {
		case <-listenChan:
			time.Sleep(10 * time.Millisecond)
			break

		case <-time.After(8 * time.Second):
			crashcounter++
			if crashcounter > 3 {
				fmt.Println("Continious crashes! Check obstructions or code -> Manual restart required")
				return
			}
			newBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "./main")
			err := newBackup.Run()
			fmt.Println("Program crashed! -> Starting new instance")
			if err != nil {
				fmt.Println("FATAL ERROR ; could not execute ./main")
				log.Fatal(err)
			}
			break

		case <-time.After(60 * time.Second):
			crashcounter = 0

		}
	}

}

func listen(listenChan chan int, udpListen *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		udpListen.ReadFromUDP(buffer[:])

		listenChan <- int(binary.LittleEndian.Uint64(buffer))
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {

	udpAddr, err := net.ResolveUDPAddr("udp", ":20014")
	if err != nil {
		log.Fatal(err)
		fmt.Println("FATAL ERROR ; could not resolve UDP address")
	}

	udpListen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("FATAL ERROR ; could not start UDP listen")
		log.Fatal(err)
	}

	newBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "./main")
	err = newBackup.Run()
	if err != nil {
		fmt.Println("FATAL ERROR ; could not execute ./main")
		log.Fatal(err)
	}

	backup(udpListen)

	udpListen.Close()

}
