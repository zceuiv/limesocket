package main

import (
	"net"
	"sync"
)

type ChannelBucket struct {
	channelMap map[int]*Channel
	mutex      *sync.Mutex
}

type Channel struct {
	id   int
	conn *net.TCPConn
}

func (cb ChannelBucket) Lock() {
	cb.mutex.Lock()
}

func (cb ChannelBucket) Unlock() {
	cb.mutex.Unlock()
}

func (cb ChannelBucket) Delete(id int) {
	cm := cb.channelMap
	if c, exist := cm[id]; exist {
		cb.Lock()
		delete(cm, id)
		c.conn.Close()
		cb.Unlock()
	}
}

func (cb ChannelBucket) Add(id int, conn *net.TCPConn) *Channel {
	c := &Channel{
		id:   id,
		conn: conn,
	}

	cb.Delete(id)

	cb.Lock()
	cb.channelMap[id] = c
	cb.Unlock()

	return c
}
