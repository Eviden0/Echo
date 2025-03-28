package test

import (
	"Ybridge/backend/cpolar"
	"encoding/json"
	"fmt"
	"testing"
)

// Tunnel 结构体，ID字段标识唯一隧道
type Tunnel struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PublicURL  string `json:"public_url"`
	Proto      string `json:"proto"`
	Addr       string `json:"addr"`
	CreateTime string `json:"create_datetime"`
}

type Response struct {
	Data struct {
		Items []struct {
			ID             string   `json:"id"` // 外层 ID 字段
			PublishTunnels []Tunnel `json:"publish_tunnels"`
		} `json:"items"`
	} `json:"data"`
}

func TestJson(t *testing.T) {
	jsonData := `{
  "data": {
    "total": 2,
    "items": [
      {
        "id": "77f4c9e2-2aca-4be6-9191-384e6dbb91e9",
        "name": "123",
        "configuration": {
          "name": "123",
          "id": "77f4c9e2-2aca-4be6-9191-384e6dbb91e9",
          "subdomain": "",
          "hostname": "",
          "proto": "http",
          "auth": "",
          "addr": "9000",
          "inspect": "false",
          "host_header": "",
          "bind_tls": "both",
          "crt": "",
          "key": "",
          "client_cas": "",
          "remote_addr": "",
          "region": "cn_top",
          "disable_keep_alives": "false",
          "redirect_https": "false",
          "start_type": "enable",
          "permanent": true
        },
        "status": "active",
        "public_url": "https://4b8715c8.r15.cpolar.top",
        "publish_tunnels": [
          {
            "name": "123",
            "public_url": "http://4b8715c8.r15.cpolar.top",
            "proto": "http",
            "addr": "http://localhost:9000",
            "type": "",
            "create_datetime": "2025-03-27T05:18:01.5916597Z"
          },
          {
            "name": "123",
            "public_url": "https://4b8715c8.r15.cpolar.top",
            "proto": "https",
            "addr": "http://localhost:9000",
            "type": "",
            "create_datetime": "2025-03-27T05:18:01.5916597Z"
          }
        ]
      },
      {
        "id": "e9d18fc0-e4a2-4a14-b121-1d0483e80ed1",
        "name": "tunnel-1",
        "configuration": {
          "name": "tunnel-1",
          "id": "e9d18fc0-e4a2-4a14-b121-1d0483e80ed1",
          "subdomain": "",
          "hostname": "",
          "proto": "http",
          "auth": "",
          "addr": "80",
          "inspect": "false",
          "host_header": "",
          "bind_tls": "both",
          "crt": "",
          "key": "",
          "client_cas": "",
          "remote_addr": "",
          "region": "cn_top",
          "disable_keep_alives": "false",
          "redirect_https": "false",
          "start_type": "enable",
          "permanent": true
        },
        "status": "active",
        "public_url": "https://1ce800d4.r15.cpolar.top",
        "publish_tunnels": [
          {
            "name": "tunnel-1",
            "public_url": "http://1ce800d4.r15.cpolar.top",
            "proto": "http",
            "addr": "http://localhost:80",
            "type": "",
            "create_datetime": "2025-03-27T05:18:09.2899782Z"
          },
          {
            "name": "tunnel-1",
            "public_url": "https://1ce800d4.r15.cpolar.top",
            "proto": "https",
            "addr": "http://localhost:80",
            "type": "",
            "create_datetime": "2025-03-27T05:18:09.2899782Z"
          }
        ]
      }
    ]
  },
  "code": 20000,
  "message": ""
}`
	// 解析 JSON
	var resp Response
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		fmt.Println("解析失败:", err)
		return
	}

	// 打印 Tunnel 信息
	for _, item := range resp.Data.Items {
		for _, tunnel := range item.PublishTunnels {
			fmt.Printf("外层ID: %s\n", item.ID)
			fmt.Printf("隧道ID: %s\n", tunnel.ID)
			fmt.Printf("名称: %s\n", tunnel.Name)
			fmt.Printf("公共URL: %s\n", tunnel.PublicURL)
			fmt.Printf("协议: %s\n", tunnel.Proto)
			fmt.Printf("地址: %s\n", tunnel.Addr)
			fmt.Printf("创建时间: %s\n", tunnel.CreateTime)
			fmt.Println("----------------------")
		}
	}
}
func TestToken(t *testing.T) {
	data := `{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMyNDM3NDIsImlhdCI6MTc0MzA3MDk0MiwiVXNlcklEIjowLCJVc2VybmFtZSI6IiIsIkVtYWlsIjoickVaS3F0SmpAZ21haWwuY29tIiwiQXBpU2VydmljZVRva2VuIjoiZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFM05ETXlORE0zTkRBc0ltbGhkQ0k2TVRjME16QTNNRGswTUN3aVZYTmxja2xFSWpvek16VXlPVEFzSWxWelpYSnVZVzFsSWpvaVprZFNVMWNpTENKRmJXRnBiQ0k2SW5KbGVtdHhkR3BxUUdkdFlXbHNMbU52YlNKOS5PZE94cmNkdUlBYmdVTVQ1NjRfMi1iTkN6UGpOM09iY0lKUmZEOGIyOGlVIn0.E99A71AuVuOZRzeUQ2nkVvsrdLTP0b6R_h8HuyqY6cc"
  },
  "code": 20000
}`
	var result map[string]interface{}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		fmt.Println("解析失败:", err)
		return
	}

	// 提取 token
	if data, ok := result["data"].(map[string]interface{}); ok {
		if token, ok := data["token"].(string); ok {
			fmt.Println(token)
		} else {
			fmt.Println("未找到 token 字段")
		}
	} else {
		fmt.Println("未找到 data 字段")
	}

}
func TestCurd(t *testing.T) {
	user := cpolar.NewUser()
	user.Login()
	user.Watch()
	user.Delete("cee9fa06-667e-433a-bdc2-e1dbd75cccb4")
	//user.Add("8999", "http")
	if len(user.Tunnels) == 0 {
		return
	}
	for _, tunnel := range user.Tunnels {
		fmt.Println(*tunnel)
	}
}
