package main

import (
	"sync"
)

type MessageBucket struct {
	msgMap map[int]chan *Message
	mutex  *sync.Mutex
}

type Message struct {
	id      int
	msgBody []byte
}

func (mb MessageBucket) Lock() {
	mb.mutex.Lock()
}

func (mb MessageBucket) Unlock() {
	mb.mutex.Unlock()
}

func (mb MessageBucket) AddChan(msg *Message) {
	id := msg.id

	var mc chan *Message
	var exist bool
	if mc, exist = mb.msgMap[id]; !exist {
		mb.Lock()
		mc = make(chan *Message, 1024)
		mb.msgMap[id] = mc
		mb.Unlock()
	}
	mc <- msg
}

func (mb MessageBucket) GetChan(id int) chan *Message {
	var mc chan *Message
	var exist bool
	if mc, exist = mb.msgMap[id]; !exist {
		mb.Lock()
		mc = make(chan *Message, 1024)
		mb.msgMap[id] = mc
		mb.Unlock()
	}
	return mc
}
