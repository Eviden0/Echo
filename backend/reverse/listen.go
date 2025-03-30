package reverse

import (
	"bufio"
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"sync"
)

const (
	tcpPort  = ":9999"
	wsPort   = ":8080"
	bufSize  = 4096
	maxConns = 100
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Session struct {
	ID          string
	TCPConn     net.Conn
	WSConn      *websocket.Conn
	FromWSChan  chan []byte
	FromTCPChan chan []byte
	CloseChan   chan struct{}
	sync.Mutex
}

var sessions sync.Map

func StartNc() {
	go initWSServer()
	initTCPServer()
}

var ID string

func GetId() string {
	return ID
}
func initWSServer() {
	http.HandleFunc("/ws", handleWSConnection)
	log.Printf("WebSocket listening on %s", wsPort)
	if err := http.ListenAndServe(wsPort, nil); err != nil {
		log.Fatalf("WebSocket server failed: %v", err)
	}
}

func initTCPServer() {
	listener, err := net.Listen("tcp", tcpPort)
	if err != nil {
		log.Fatalf("TCP listen failed: %v", err)
	}
	defer listener.Close()
	log.Printf("TCP listening on %s", tcpPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept failed: %v", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(tcpConn net.Conn) {
	defer tcpConn.Close()
	sessionID := tcpConn.RemoteAddr().String()
	ID = sessionID
	session := &Session{
		ID:          sessionID,
		TCPConn:     tcpConn,
		FromWSChan:  make(chan []byte, bufSize),
		FromTCPChan: make(chan []byte, bufSize),
		CloseChan:   make(chan struct{}),
	}
	sessions.Store(sessionID, session)

	log.Printf("New TCP session: %s", sessionID)

	tcpConn.Write([]byte("/usr/bin/script -qc /bin/bash /dev/null\n"))

	var wg sync.WaitGroup
	wg.Add(2)

	// 从TCP读取数据转发到WebSocket
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(tcpConn)
		for {
			select {
			case <-session.CloseChan:
				return
			default:
				data := make([]byte, bufSize)
				n, err := reader.Read(data)
				if err != nil {
					close(session.CloseChan)
					return
				}
				session.FromTCPChan <- bytes.TrimSpace(data[:n])
			}
		}
	}()

	// 从WebSocket通道写数据到TCP
	go func() {
		defer wg.Done()
		for {
			select {
			case data := <-session.FromWSChan:
				if _, err := tcpConn.Write(append(data, '\n')); err != nil {
					close(session.CloseChan)
					return
				}
			case <-session.CloseChan:
				return
			}
		}
	}()

	wg.Wait()
	sessions.Delete(sessionID)
	log.Printf("Session terminated: %s", sessionID)
}

func handleWSConnection(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer wsConn.Close()

	sessionID := r.URL.Query().Get("session") //对应相应的id
	/*
		后面改成将这一串发送给前端然后再启动前端
	*/
	if sessionID == "" {
		log.Println("Missing session ID")
		return
	}

	val, ok := sessions.Load(sessionID)
	if !ok {
		log.Printf("Invalid session ID: %s", sessionID)
		return
	}
	session := val.(*Session)
	session.Lock()
	session.WSConn = wsConn
	session.Unlock()

	log.Printf("WebSocket connected to session: %s", sessionID)

	var wg sync.WaitGroup
	wg.Add(2)

	// 从WebSocket读取指令
	go func() {
		defer wg.Done()
		for {
			_, msg, err := wsConn.ReadMessage()
			if err != nil {
				//close(session.CloseChan)
				return
			}
			session.FromWSChan <- bytes.TrimSpace(msg)
		}
	}()

	// 向WebSocket发送TCP响应
	go func() {
		defer wg.Done()
		for {
			select {
			case data := <-session.FromTCPChan:
				if err := wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
					close(session.CloseChan)
					return
				}
			case <-session.CloseChan:
				return
			}
		}
	}()
	wg.Wait()
}

//func (s *Session) Stop() {
//	s.Lock()
//	defer s.Unlock()
//
//	// Close the TCP connection
//	if s.TCPConn != nil {
//		s.TCPConn.Close()
//		s.TCPConn = nil
//	}
//
//	// Close the WebSocket connection
//	if s.WSConn != nil {
//		s.WSConn.Close()
//		s.WSConn = nil
//	}
//
//	// Close all channels
//	close(s.CloseChan)
//	close(s.FromWSChan)
//	close(s.FromTCPChan)
//
//	// Reset channels
//	s.FromWSChan = make(chan []byte, bufSize)
//	s.FromTCPChan = make(chan []byte, bufSize)
//	s.CloseChan = make(chan struct{})
//
//	log.Printf("Session stopped: %s", s.ID)
//}
