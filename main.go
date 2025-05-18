package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	reuse "github.com/libp2p/go-reuseport"
)

var (
	localaddr  = "0.0.0.0:7684"
	serveraddr = "112.125.89.8:45914"
	remoteaddr = "1.1.1.1:111"

	role = 1 //1:peer, 2:host
)

func prompt(pm string) string {
	fmt.Print(pm)
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	return string(data)
}

func tryPunch(dietun chan int, conntun chan net.Conn) {
	for {
		select {
		case <-dietun:
			return
		default:
			time.Sleep(500 * time.Millisecond)
			peerConn, err := reuse.DialTimeout("tcp", localaddr, remoteaddr, 2*time.Second)
			// fmt.Println(peeraddr)
			// peerConn, err := reuse.Dial("tcp", localaddr, peeraddr)
			if err != nil {
				fmt.Println("Peer Conn Error:", err)
			} else {
				select {
				case <-dietun:
					peerConn.Close()
					return
				case conntun <- peerConn:
					fmt.Println("Peer Connected")
					return
				}
			}
		}
	}
}

func doListenPeer(dietun chan int, conntun chan net.Conn) { //host only
	listener, _ := reuse.Listen("tcp", localaddr)

	select {
	case <-dietun:
		return
	default:
		peerConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept Error:", err)
		}

		select {
		case conntun <- peerConn:
			fmt.Println("Accepted Peer Connect")
			return
		case <-dietun:
			peerConn.Close()
			return
		}
	}
}

func main() {
	fmt.Println("TCP Punching Hole & Forwarder Test")

	serverConn, err := reuse.Dial("tcp", localaddr, serveraddr)
	if err != nil {
		fmt.Println("Error:", err)
	}

	remoteaddr, _ = bufio.NewReader(serverConn).ReadString('\n')
	remoteaddr = strings.TrimSpace(remoteaddr)
	fmt.Println("Remote addr:", remoteaddr)

	dietun := make(chan int, 1)
	conntun := make(chan net.Conn, 1)

	go tryPunch(dietun, conntun)
	go doListenPeer(dietun, conntun)

	peerConn := <-conntun
	remoteconn = peerConn

	if role == 1 { //peer
		go handleRemotePeer(peerConn)
		go ListenLocal()
	} else { //host
		go handleRemoteHost(peerConn)
	}

	// useless recv
	// go func() {
	// 	for {
	// 		buf := make([]byte, 1024)
	// 		n, err := peerConn.Read(buf)
	// 		if err != nil {
	// 			fmt.Println("Recv Error:", err)
	// 		}

	// 		fmt.Println("Recv From Peer:", string(buf[:n]))
	// 	}
	// }()

shell:
	for {
		input := prompt(">")
		args := strings.SplitN(input, " ", 2)

		switch args[0] {
		case "server":
			if len(args) < 2 {
				fmt.Println(len(args))
				break
			}
			_, err := peerConn.Write([]byte(args[1]))
			if err != nil {
				fmt.Println(err)
			}
		case "peer":
			if len(args) < 2 {
				break
			}
			_, err := peerConn.Write([]byte(args[1]))
			if err != nil {
				fmt.Println(err)
			}
		case "exit":
			fmt.Println("Bye!")
			break shell
		}
	}
}
