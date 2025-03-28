package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"sync"
	"testing"
)

func TestNc(t *testing.T) {
	listen, err := net.Listen("tcp", ":9999")
	log.Println("listen on :9999")
	if err != nil {
		t.Logf("%v 开启监听失败", err)
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			t.Logf("%v 接收连接失败", err)
		}
		go handleConn2(conn)
	}
}
func handleConn2(conn net.Conn) {
	var wg sync.WaitGroup
	localChan := make(chan []byte)
	remoteChan := make(chan []byte)

	// 协程1：读取本地输入到remoteChan
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(os.Stdin)
		for {
			input, err := reader.ReadBytes('\n')
			if err != nil {
				close(remoteChan)
				return
			}
			remoteChan <- input
		}
	}()

	// 协程2：读取远程连接数据到localChan
	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				close(localChan)
				return
			}
			localChan <- buf[:n]
		}
	}()

	// 协程3：发送本地指令到远程连接
	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range remoteChan {
			if _, err := conn.Write(data); err != nil {
				return
			}
		}
	}()

	// 协程4：输出远程返回到本地终端
	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range localChan {
			os.Stdout.Write(data)
		}
	}()

	wg.Wait()
	conn.Close()
}
