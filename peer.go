package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

var localforaddr = "127.0.0.1:17789"

func handleRemotePeer(conn net.Conn) {
	log.Println("Remote:", conn.RemoteAddr().String())
	for {
		idbuf := make([]byte, 4)
		_, err := io.ReadFull(conn, idbuf)
		if err != nil {
			return
		}
		id := binary.BigEndian.Uint32(idbuf)
		log.Println("0id", id)

		peerconn := getConnp(id)
		if peerconn == nil {
			log.Fatalln("[RECVBADID]")
			continue
		}

		lenbuf := make([]byte, 4)
		io.ReadFull(conn, lenbuf)
		len := binary.BigEndian.Uint32(lenbuf)
		log.Println("0len", len)

		data := make([]byte, len)
		io.ReadFull(conn, data)

		_, err = peerconn.Write(data)
		if err != nil {
			log.Fatalln("[RECVERR]", err)
			delConnp(id)
		}
	}
}

func handleLocalPeer(conn net.Conn, id uint32) {
	log.Println("LOCAL:", conn.RemoteAddr().String())
	for {
		data := []byte{}
		buf := make([]byte, 1024)

		n, err := conn.Read(buf)
		if err != nil {
			delConnp(id)
			return
		}

		idbyte := make([]byte, 4)
		binary.BigEndian.PutUint32(idbyte, id)
		log.Println("1idb", idbyte)
		// _, err = remoteconn.Write(idbyte)
		// if err != nil {
		// 	log.Fatalln("[WRITEERR]", err)
		// 	return
		// }
		data = append(data, idbyte...)

		len := make([]byte, 4)
		binary.BigEndian.PutUint32(len, uint32(n))
		log.Println("1lenb", len)
		// _, err = remoteconn.Write(len)
		// if err != nil {
		// 	log.Fatalln("[WRITEERR]", err)
		// 	return
		// }
		data = append(data, len...)

		data = append(data, buf[:n]...)
		_, err = remoteconn.Write(data)
		// _, err = remoteconn.Write(buf[:n])
		if err != nil {
			log.Fatalln("[WRITEERR]", err)
			return
		}
	}
}

func ListenLocal() {
	listenerf, err := net.Listen("tcp", localforaddr)
	if err != nil {
		log.Fatalln("[ERROR]", err)
	}

	for {
		conn, err := listenerf.Accept()
		if err != nil {
			log.Fatalln("[ERROR]", err)
		}

		id := addConnp(conn)

		go handleLocalPeer(conn, id)
	}
}
