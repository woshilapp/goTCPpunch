package main

import (
	"net"
	"sync"
)

var ID uint32 = 0
var connlist map[uint32]net.Conn = map[uint32]net.Conn{}
var cllock sync.RWMutex
var remoteconn net.Conn

func addConnp(conn net.Conn) uint32 {
	cllock.Lock()
	defer cllock.Unlock()

	ID++
	id := ID - 1

	connlist[id] = conn

	return id
}

func delConnp(id uint32) {
	cllock.Lock()
	defer cllock.Unlock()

	if _, exist := connlist[id]; !exist {
		return
	}

	delete(connlist, id)
}

func getConnp(id uint32) net.Conn {
	cllock.RLock()
	defer cllock.RUnlock()

	if _, exist := connlist[id]; !exist {
		return nil
	}

	return connlist[id]
}

func addConnh(id uint32, conn net.Conn) {
	cllock.Lock()
	defer cllock.Unlock()

	connlist[id] = conn
}

func delConnh(id uint32) {
	cllock.Lock()
	defer cllock.Unlock()

	if _, exist := connlist[id]; !exist {
		return
	}

	delete(connlist, id)
}

func getConnh(id uint32) net.Conn {
	cllock.RLock()
	defer cllock.RUnlock()

	if _, exist := connlist[id]; !exist {
		return nil
	}

	return connlist[id]
}
