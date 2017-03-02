package main

import (
	log "github.com/alecthomas/log4go"
	"sync"
)

var (
	OnlineMsg   chan *Message
	OfflineMsg  MessageBucket
	UserChannel ChannelBucket
)

func init() {
	err := InitConfig()
	if err != nil {
		log.Error("Init config error: %v", err)
	}
	// init log
	log.LoadConfiguration(Conf.Log)

	UserChannel = ChannelBucket{
		channelMap: make(map[int]*Channel),
		mutex:      &sync.Mutex{},
	}
	OfflineMsg = MessageBucket{
		msgMap: make(map[int]chan *Message),
		mutex:  &sync.Mutex{},
	}
	OnlineMsg = make(chan *Message, 1024)
}

func main() {

	defer log.Close()

	StartTCP()
	StartHTTP()
	signalCH := InitSignal()
	HandleSignal(signalCH)
	log.Info("server stop")
}
