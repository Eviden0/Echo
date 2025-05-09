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
	tcpPort = ":9999"
	wsPort  = ":8080"
	bufSize = 4096
)

// 升级 http 到 websocket 协议
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
	http.HandleFunc("/close", handleCloseRequest) // 可选：通过 HTTP 关闭 session
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
	tcpConn.Write([]byte("clear\n"))

	tcpConn.Write([]byte("/usr/bin/script -qc /bin/bash /dev/null\n"))

	var wg sync.WaitGroup
	wg.Add(2)

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
					//  优化了关闭conn的逻辑,fixed panic close closed chan
					select {
					case <-session.CloseChan:
					default:
						close(session.CloseChan)
					}
					return
				}
				session.FromTCPChan <- bytes.TrimSpace(data[:n])
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case data := <-session.FromWSChan:
				if _, err := tcpConn.Write(append(data, '\n')); err != nil {
					select {
					case <-session.CloseChan:
					default:
						close(session.CloseChan)
					}
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

	sessionID := r.URL.Query().Get("session")
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

	go func() {
		defer wg.Done()
		for {
			_, msg, err := wsConn.ReadMessage()
			if err != nil {
				return
			}
			session.FromWSChan <- bytes.TrimSpace(msg)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case data := <-session.FromTCPChan:
				if err := wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
					select {
					case <-session.CloseChan:
					default:
						close(session.CloseChan)
					}
					return
				}
			case <-session.CloseChan:
				return
			}
		}
	}()

	wg.Wait()
}

// 防止重复关闭资源
func closeSession(session *Session) {
	session.Lock()
	defer session.Unlock()

	select {
	case <-session.CloseChan:
		return
	default:
		close(session.CloseChan)
	}

	if session.TCPConn != nil {
		_ = session.TCPConn.Close()
		session.TCPConn = nil
	}

	if session.WSConn != nil {
		_ = session.WSConn.Close()
		session.WSConn = nil
	}

	sessions.Delete(session.ID)
	log.Printf("Session %s closed", session.ID)
}

// ✅ 外部通过 ID 主动关闭
func CloseConnByID(sessionID string) {
	val, ok := sessions.Load(sessionID)
	if !ok {
		log.Printf("[GracefulExit] No session found with ID %s", sessionID)
		return
	}
	session := val.(*Session)
	log.Printf("[GracefulExit] Gracefully closing session %s", sessionID)
	closeSession(session)
}

// 提供一个restful接口来关闭session,eg: http://localhost:8080/close?session=SessionID
func handleCloseRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源

	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		http.Error(w, "Missing session ID", http.StatusBadRequest)
		return
	}
	CloseConnByID(sessionID)
	w.Write([]byte("Session closed: " + sessionID))
}
