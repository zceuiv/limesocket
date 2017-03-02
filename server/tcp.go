package main

import (
	"fmt"
	log "github.com/alecthomas/log4go"
	"net"
)

func StartTCP() {
	// 这里暂时只监听一个地址端口，不然发消息协程不支持多地址端口
	for _, bind := range Conf.TCPBind {
		log.Info("start tcp listen addr:\"%s\"", bind)
		go tcpListen(bind)
	}

}

func tcpListen(bind string) {
	addr, err := net.ResolveTCPAddr("tcp", bind)
	if err != nil {
		log.Error("net.ResolveTCPAddr(\"tcp\"), %s) error(%v)", bind, err)
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Error("net.ListenTCP(\"tcp4\", \"%s\") error(%v)", bind, err)
		panic(err)
	}
	// free the listener resource
	defer func() {
		log.Info("tcp addr: \"%s\" close", bind)
		if err := listener.Close(); err != nil {
			log.Error("listener.Close() error(%v)", err)
		}
	}()

	go onlineMsgRoutine()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Error("listener.AcceptTCP() error(%v)", err)
			continue
		}
		if err = conn.SetKeepAlive(Conf.TCPKeepalive); err != nil {
			log.Error("conn.SetKeepAlive() error(%v)", err)
			conn.Close()
			continue
		}
		if err = conn.SetReadBuffer(Conf.RcvbufSize); err != nil {
			log.Error("conn.SetReadBuffer(%d) error(%v)", Conf.RcvbufSize, err)
			conn.Close()
			continue
		}
		if err = conn.SetWriteBuffer(Conf.SndbufSize); err != nil {
			log.Error("conn.SetWriteBuffer(%d) error(%v)", Conf.SndbufSize, err)
			conn.Close()
			continue
		}

		var c *Channel
		if c, err = validateConn(conn); err != nil {
			log.Error("validateConn(%v) error(%v)", *conn, err)
			conn.Close()
			continue
		}

		// one connection one routine
		go handleTCPConn(c)
		log.Debug("accept finished")
	}
}

func validateConn(conn *net.TCPConn) (*Channel, error) {
	tokenP := make([]byte, 256)
	tokenLen, err := conn.Read(tokenP)
	if err != nil {
		return nil, err
	}
	token := string(tokenP[0:tokenLen])
	id, err := validateToken(token)
	if err != nil {
		return nil, err
	}

	c := UserChannel.Add(id, conn)

	fmt.Println("validated:", id)
	return c, nil
}

func handleTCPConn(c *Channel) {
	id := c.id
	conn := c.conn
	addr := conn.RemoteAddr().String()
	log.Debug("<%s> handleTcpConn routine start", addr)

	sendStoredMsg(id)

	p := make([]byte, 1024)
	for {
		len, err := conn.Read(p)
		if err != nil {
			UserChannel.Delete(id)
			fmt.Printf("connection read error: %v\n", err)
			break
		}
		handleRcvMsg(id, p[0:len])
	}
	return
}

func handleRcvMsg(id int, rcvMsg []byte) {
	fmt.Println(id, string(rcvMsg))

	resp := []byte("Server has received message: ")
	for _, p := range rcvMsg {
		resp = append(resp, p)
	}
	fmt.Println("111111111111", string(resp))
	msg := &Message{
		id:      id,
		msgBody: resp,
	}

	OnlineMsg <- msg
}

func onlineMsgRoutine() {
	for {
		msg := <-OnlineMsg
		sendMsg(msg)
	}
}

func storeMsg(msg *Message) {
	OfflineMsg.AddChan(msg)
}

//第一个返回值表示是否在线
func sendMsg(msg *Message) (bool, error) {
	id := msg.id

	var c *Channel
	var exist bool
	if c, exist = UserChannel.channelMap[id]; !exist {
		storeMsg(msg)
		return false, nil
	}

	_, err := c.conn.Write(msg.msgBody)
	if err != nil {
		return false, err
	}
	return true, nil
}

func sendStoredMsg(id int) {
	ch := OfflineMsg.GetChan(id)
	for {
		select {
		case msg := <-ch:
			fmt.Println("send stored message")
			ret, err := sendMsg(msg)
			if err != nil {
				log.Error("sendMsg(%v) error(%v)", *msg, err)
				return
			} else if !ret {
				log.Debug("sendMsg(%v) offline", *msg)
				return
			} else {
				continue
			}
		default:
			return
		}
	}
}
