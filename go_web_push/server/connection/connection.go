package connection

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	wsConn *websocket.Conn
	inChan chan []byte
	outChan chan []byte
	closeChan chan byte
	mutex sync.Mutex
	isClose bool
}


func InitConnection (wsConn *websocket.Conn)(conn *Connection,err error){
	conn = &Connection{
		wsConn:wsConn,
		inChan:make(chan []byte,1000),
		outChan:make(chan []byte,1000),
		closeChan:make(chan byte,1),
	}
	//启动读写协程
	go conn.readLopp()
	go conn.writeLopp()
	return
}

func (conn *Connection)ReadMessag() (data []byte,err error){
	select{
	case data=<-conn.inChan:
	case <-conn.closeChan:
		err=errors.New("connection is closed")
	}
	return
}


func(conn *Connection)WriteMessag(data []byte)(err error){
	select{
	case conn.outChan<-data:
	case <-conn.closeChan:
		err=errors.New("connection is closed")
	}
	return
}

func (conn *Connection)Close(){
	//线程安全。可并发
	conn.wsConn.Close()

	//这里的代码只执行一次 加锁
	conn.mutex.Lock()
	if !conn.isClose {
		close(conn.closeChan)
		conn.isClose = true
	}
	conn.mutex.Unlock()
}

//内部实现消息

func (conn *Connection)readLopp(){
	var(
		data []byte
		err error
	)
	for {
		if _,data,err=conn.wsConn.ReadMessage();err !=nil{
			goto ERR
		}
		//这里会堵塞
		select{
		case conn.inChan<-data:
		case<-conn.closeChan:
			goto ERR
		}
	}
ERR:
	conn.Close()
}

func(conn *Connection)writeLopp(){
	var(
		data []byte
		err error
	)
	for {
		select{
		case data=<-conn.outChan:
		case <-conn.closeChan:
			goto ERR
		}
		if err = conn.wsConn.WriteMessage(websocket.TextMessage,data);err!=nil{
			goto ERR

		}
	}
ERR:
	conn.Close()
}
