package main

import (
	"flag"
	"runtime"
)

var (
	Conf     *Config
	confFile string
)

func init() {
	flag.StringVar(&confFile, "c", "./comet.conf", " set gopush-cluster comet config file path")
}

type Config struct {
	// base
	Log      string   `goconf:"base:log"`
	HTTPBind []string `goconf:"base:http.bind:,"`
	TCPBind  []string `goconf:"base:tcp.bind:,"`
	// channel
	SndbufSize    int  `goconf:"channel:sndbuf.size:memory"`
	RcvbufSize    int  `goconf:"channel:rcvbuf.size:memory"`
	BufioInstance int  `goconf:"channel:bufio.instance"`
	BufioNum      int  `goconf:"channel:bufio.num"`
	TCPKeepalive  bool `goconf:"channel:tcp.keepalive"`
}

// InitConfig get a new Config struct.
func InitConfig() error {
	Conf = &Config{
		// base
		Log:      "./log.xml",
		HTTPBind: []string{"127.0.0.1:10079"},
		TCPBind:  []string{"127.0.0.1:10080"},
		// channel
		SndbufSize:    2048,
		RcvbufSize:    256,
		BufioInstance: runtime.NumCPU(),
		BufioNum:      128,
		TCPKeepalive:  false,
	}
	return nil
}
