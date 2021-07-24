package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go_Multi_learn/go_web_push/server/connection"
	"html/template"
	"net/http"
	"time"
)

var(
	upgrader = websocket.Upgrader{
		//握手的过程中跨域问题 要允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},

	}
)
type MyHttpHandler struct{}
func wsHandler(res http.ResponseWriter,rsp *http.Request){
	var(
		wsConn *websocket.Conn
		err error
		conn *connection.Connection
		data []byte
	 )
	//完成ws协议和http协议的转换
    if wsConn,err =upgrader.Upgrade(res,rsp,nil);err!=nil{
    	fmt.Println(err)
		return
	}
    fmt.Println(res,rsp)
    if conn,err=connection.InitConnection(wsConn);err!=nil{
		goto ERR
	}
//给前端定义的心跳
    go func() {
    	var(
    		err error
		)
    	for {
			if err=conn.WriteMessag([]byte("heartbeat"));err!=nil{
				goto ERR
			}
			time.Sleep(time.Second*1)
		}
    ERR:
	conn.Close()
	}()

    for{
    	if data,err =conn.ReadMessag();err!=nil{
    		goto ERR
		}
    	if err=conn.WriteMessag(data);err!=nil{
    		goto ERR
		}
	}
ERR:
	conn.Close()

}

func indexHandler(w http.ResponseWriter,rsp *http.Request){
	t:=template.Must(template.ParseFiles("./go_web_push/client.html"))
	t.Execute(w, time.Now().Format("2006-01-02 15:04:05"))
}
func main(){

	http.HandleFunc("/ws",wsHandler)
	http.HandleFunc("/index",indexHandler)
	http.ListenAndServe(":8088", nil)
}


