package main

import (
	"encoding/json"
	"fmt"
	log "github.com/alecthomas/log4go"
	"io"
	"io/ioutil"
	"net/http"
)

func StartHTTP() {

	for _, bind := range Conf.HTTPBind {
		log.Info("start http listen addr:\"%s\"", bind)
		go httpListen(bind)
	}
}

func httpListen(bind string) {
	var err error

	mux := &http.ServeMux{}

	mux.HandleFunc("/testpub/", func(w http.ResponseWriter, r *http.Request) {
		hasBody := (r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH") && r.Body != nil
		ct := r.Header.Get("Content-Type")

		if hasBody && ct != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: Content-Type must be application/json.\n")))
			return
		}

		var bodyJson map[string]interface{}
		bodyJson, err = readBodyJson(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: %v\n", err)))
			return
		}

		var ok bool = true

		var id float64
		if id, ok = bodyJson["id"].(float64); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("id must be a number."))
			return
		}

		var msgbody string
		if msgbody, ok = bodyJson["msg"].(string); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("msg must be a string."))
			return
		}

		msg := &Message{id: int(id), msgBody: []byte(msgbody)}
		OnlineMsg <- msg

		w.Write([]byte(fmt.Sprintf("Success: %v\n", *msg)))
	})

	http.ListenAndServe(bind, mux)
}

func readBodyJson(r io.Reader) (bodyJson map[string]interface{}, err error) {
	var b []byte
	b, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &bodyJson)
	return
}
