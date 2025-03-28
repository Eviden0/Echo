package test

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"testing"
)

func TestEnv(t *testing.T) {
	// 获取 %USERPROFILE% 环境变量的值
	confiPath := os.Getenv("USERPROFILE") + "\\.cpolar\\cpolar.yml"
	file, err := os.OpenFile(confiPath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	content, _ := io.ReadAll(file)
	t.Log(string(content))
}
func TestIP(t *testing.T) {
	addr := "tcp://15.tcp.cpolar.top:13375"
	// Extract the domain part from the address
	parts := strings.Split(addr, "://")

	// Lookup the IP address for the domain
	domain := strings.Split(parts[1], ":")[0]
	ips, _ := net.LookupIP(domain)

	// Replace the domain part with the IP address
	ip := ips[0].String()
	port := strings.Split(parts[1], ":")[1]
	fmt.Println(parts[0] + "://" + ip + ":" + port)
}
