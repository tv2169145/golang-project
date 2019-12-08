package main

import (
	"fmt"
	"github.com/tv2169145/golang-project/grpc/04-chat/chat"
	"google.golang.org/grpc"
	"io"
	"net"
	"sync"
)

// 封裝單一連線
type Connection struct {
	conn chat.Chat_ChatServer // 單一server
	send chan *chat.ChatMessage // 回傳的message channel
	quit chan struct{} // 判斷是否退出的channel
}

// 初始化並封裝單一連線, 同時開始監聽channel
func NewConnection(conn chat.Chat_ChatServer) *Connection {
	c := &Connection{
		conn: conn,
		send: make(chan *chat.ChatMessage),
		quit: make(chan struct{}),
	}
	go c.start()
	return c
}

// 啟動監聽是否有訊息傳入send通道, 並將訊息傳給user
func (c *Connection) start() {
	running := true
	for running {
		select {
		case msg := <-c.send:
			c.conn.Send(msg) // Ignoring the error, they just don't get this message.
		case <-c.quit:
			running = false
		}
	}
}

// 將 chatMessage 丟入待發送的channel內
func (c *Connection) Send(msg *chat.ChatMessage) {
	defer func() {
		// Ignore any errors about sending on a closed channel
		recover()
	}()
	c.send <- msg
}

// 關閉/中斷單一連線
func (c *Connection) Close() error {
	close(c.send)
	close(c.quit)
	return nil
}

// 取得user輸入的chatMessage, 並將message丟入推播池
func (c *Connection) GetMessages(broadcast chan<- *chat.ChatMessage) error {
	for {
		msg, err := c.conn.Recv()
		if err == io.EOF {
			c.Close()
			return nil
		} else if err != nil {
			c.Close()
			return err
		}
		go func(msg *chat.ChatMessage) {
			select {
			case broadcast<- msg:
			case <-c.quit:
			}
		}(msg)
	}
}

// 封裝多個連線, 建立連線池-----------------------------------------------------------------
type ChatServer struct {
	broadcast chan *chat.ChatMessage	// 推播池, 將所有連線的訊息集中在此,準備將其回傳給user
	quit chan struct{}	// 關閉所有連線
	connections []*Connection	// 	收集所有連線(連線池)
	connLock sync.Mutex	// 併發鎖
}

// 初始化連線池, 並開始監聽各連線
func NewChatServer() *ChatServer {
	srv := &ChatServer{
		broadcast: make(chan *chat.ChatMessage),
		quit: make(chan struct{}),
	}
	go srv.start()
	return srv
}

func (c *ChatServer) Close() error {
	close(c.quit)
	return nil
}

// 啟動監聽連線池內所有連線, 並將推播池內訊息發送給所有連線者
func (c *ChatServer) start() {
	running := true
	for running {
		select {
		case msg := <-c.broadcast:
			c.connLock.Lock()
			for _, v := range c.connections {
				go v.Send(msg)
			}
			c.connLock.Unlock()
		case <-c.quit:
			running = false
		}
	}
}

func (c *ChatServer) Chat(stream chat.Chat_ChatServer) error {
	// 建立一個新的連線
	conn := NewConnection(stream)

	// 將此連線加入連線池
	c.connLock.Lock()
	c.connections = append(c.connections, conn)
	c.connLock.Unlock()

	// 開始監聽此連線, 將收到的訊息加入推播池
	err := conn.GetMessages(c.broadcast)

	c.connLock.Lock()
	for i, v := range c.connections {
		if v == conn {
			c.connections = append(c.connections[:i], c.connections[i+1:]...)
		}
	}
	c.connLock.Unlock()
	return err
}



func main() {
	lst, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	srv := NewChatServer()
	chat.RegisterChatServer(s, srv)
	fmt.Println("Now serving at port 8080")
	err = s.Serve(lst)
	if err != nil {
		panic(err)
	}
}
