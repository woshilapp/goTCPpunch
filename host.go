package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

var localseraddr = "127.0.0.1:6699"

func handleRemoteHost(conn net.Conn) {
	log.Println("Remote:", conn.RemoteAddr().String())
	for {
		idbuf := make([]byte, 4)
		_, err := io.ReadFull(conn, idbuf)
		if err != nil {
			return
		}
		id := binary.BigEndian.Uint32(idbuf)
		log.Println("0id", id)

		peerconn := getConnh(id)
		if peerconn == nil {
			peerconn, err := net.Dial("tcp", localseraddr)
			if err != nil {
				log.Fatalln("[RECVERR]", err)
				continue
			}

			addConnh(id, peerconn)

			go handleLocalHost(peerconn, id)
		}

		peerconn = getConnh(id)

		lenbuf := make([]byte, 4)
		io.ReadFull(conn, lenbuf)
		len := binary.BigEndian.Uint32(lenbuf)
		log.Println("0len", len)

		data := make([]byte, len)
		io.ReadFull(conn, data)

		_, err = peerconn.Write(data)
		if err != nil {
			log.Fatalln("[RECVERR]", err)
			delConnh(id)
		}
	}
}

func handleLocalHost(conn net.Conn, id uint32) {
	log.Println("LOCAL:", conn.RemoteAddr().String())
	for {
		data := []byte{}
		buf := make([]byte, 1024)

		n, err := conn.Read(buf)
		if err != nil {
			delConnh(id)
			return
		}

		idbyte := make([]byte, 4)
		binary.BigEndian.PutUint32(idbyte, id)
		log.Println("1idb", idbyte)
		data = append(data, idbyte...)
		// _, err = remoteconn.Write(idbyte)
		// if err != nil {
		// 	log.Fatalln("[WRITEERR]", err)
		// 	return
		// }

		len := make([]byte, 4)
		binary.BigEndian.PutUint32(len, uint32(n))
		log.Println("1lenb", len)
		data = append(data, len...)
		// _, err = remoteconn.Write(len)
		// if err != nil {
		// 	log.Fatalln("[WRITEERR]", err)
		// 	return
		// }

		data = append(data, buf[:n]...)
		_, err = remoteconn.Write(data)
		// _, err = remoteconn.Write(buf[:n])
		if err != nil {
			log.Fatalln("[WRITEERR]", err)
			return
		}
	}
}
