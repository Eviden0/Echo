package cpolar

import (
	"Ybridge/backend/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type User struct {
	email    string
	password string
	cookie   string
	data     interface{} //拿到一个信息或者错误
	client   *http.Client
	Tunnels  []*Tunnel
}
type Tunnel struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PublicURL  string `json:"public_url"`
	Proto      string `json:"proto"`
	LocalAddr  string `json:"addr"` //本地地址
	CreateTime string `json:"create_datetime"`
}

func NewUser() *User {
	var client *http.Client
	if config.IsDebug {
		proxyURL, _ := url.Parse("http://localhost:8083/")
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
				// 其他 Transport 配置（可选）
				// DialContext: (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
			},
		}
	} else {
		client = &http.Client{}
	}
	return &User{
		email:    config.Email,
		password: config.Password,
		client:   client,
	}
}
func Run() *User {
	user := NewUser()
	user.Login()
	user.Watch()
	return user
}

// Mode
const (
	Login = iota
	Add
	Delete
	Watch
)

// Login Mode 0
func (u *User) Login() {
	data := map[string]interface{}{
		"email":    u.email,
		"password": u.password,
	}
	jsonData, _ := json.Marshal(data)
	// Login 更新cookie字段
	req, _ := http.NewRequest("POST", "http://localhost:9200/api/v1/user/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", " application/json;charset=utf-8")
	resp, _ := u.client.Do(req)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("解析失败:", err)
		return
	}
	// 提取 token
	if data, ok := result["data"].(map[string]interface{}); ok {
		if token, ok := data["token"].(string); ok {
			u.cookie = token
		} else {
			fmt.Println("未找到 token 字段")
		}
	} else {
		fmt.Println("未找到 data 字段")
	}
}

// Add Mode 1
func (u *User) Add(port, proto string) {
	//	 确认协议和端口
	data := fmt.Sprintf(`{"name":"TunByEcho","proto":"%v","addr":"%v","subdomain":"","hostname":"","auth":"","inspect":"false","host_header":"","bind_tls":"both","remote_addr":"","region":"cn_top","disable_keep_alives":"false","redirect_https":"false","start_type":"enable","permanent":true,"crt":"","key":"","client_cas":""}
`, proto, port)
	req, _ := http.NewRequest("POST", "http://localhost:9200/api/v1/tunnels", bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", " application/json;charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+u.cookie)
	_, _ = u.client.Do(req)
}

// Delete Mode 2
func (u *User) Delete(id string) {
	req, _ := http.NewRequest("DELETE", "http://localhost:9200/api/v1/tunnels/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+u.cookie)
	resp, _ := u.client.Do(req)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if strings.Contains(string(body), "\"code\": 20000,") {
		log.Printf("删除成功 %s\n", id)
	} else {
		log.Printf("删除失败 %s\n", id)
	}
	//更新本地Tunnel信息
	for i, tunnel := range u.Tunnels {
		if tunnel.ID == id {
			if len(u.Tunnels) > 1 {
				// 删除两个
				if tunnel.Proto == "http" {
					u.Tunnels = append(u.Tunnels[:i], u.Tunnels[i+2:]...)
				} else {
					u.Tunnels = append(u.Tunnels[:i-1], u.Tunnels[i+1:]...)
				}

			} else {
				u.Tunnels = []*Tunnel{}
			}
		}
	}
}
func getIP(addr string) string {
	// Check if the protocol is TCP
	if !strings.HasPrefix(addr, "tcp://") {
		return addr // Return the original address if it's not TCP
	}

	// Extract the domain part from the address
	parts := strings.Split(addr, "://")
	if len(parts) != 2 {
		return addr // Return the original address if it doesn't match the expected format
	}

	// Lookup the IP address for the domain
	domain := strings.Split(parts[1], ":")[0]
	ips, err := net.LookupIP(domain)
	if err != nil || len(ips) == 0 {
		return addr // Return the original address if the lookup fails
	}

	// Replace the domain part with the IP address
	ip := ips[0].String()
	port := strings.Split(parts[1], ":")[1]
	return parts[0] + "://" + ip + ":" + port
}

// Watch 的同时及时更新Tunnel信息,Mode 3
func (u *User) Watch() {
	type jsonResponse struct {
		Data struct {
			Items []struct {
				ID             string   `json:"id"` // 外层 ID 字段
				PublishTunnels []Tunnel `json:"publish_tunnels"`
			} `json:"items"`
		} `json:"data"`
	}
	req, _ := http.NewRequest("GET", "http://localhost:9200/api/v1/tunnels", nil)
	req.Header.Set("Authorization", "Bearer "+u.cookie)
	resp, _ := u.client.Do(req)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var jsonresp jsonResponse
	if err := json.Unmarshal(body, &jsonresp); err != nil {
		fmt.Println("解析错误:", err)
	}
	// 打印 Tunnel 信息
	for _, item := range jsonresp.Data.Items {
		for _, tunnel := range item.PublishTunnels {
			t := &Tunnel{
				ID:         item.ID, // 外层 ID
				Name:       tunnel.Name,
				PublicURL:  getIP(tunnel.PublicURL), //添加ip+port
				Proto:      tunnel.Proto,
				LocalAddr:  tunnel.LocalAddr,
				CreateTime: tunnel.CreateTime,
			}

			// 将 Tunnel 添加到 User 的 Tunnels 切片中
			u.Tunnels = append(u.Tunnels, t)
			//fmt.Printf("隧道ID: %s\n", item.ID)
			//fmt.Printf("名称: %s\n", tunnel.Name)
			//fmt.Printf("公共URL: %s\n", tunnel.PublicURL)
			//fmt.Printf("协议: %s\n", tunnel.Proto)
			//fmt.Printf("地址: %s\n", tunnel.Addr)
			//fmt.Printf("创建时间: %s\n", tunnel.CreateTime)
			//fmt.Println("----------------------")
		}
	}
}
