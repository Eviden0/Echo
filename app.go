package main

import (
	"Ybridge/backend/config"
	"Ybridge/backend/cpolar"
	"Ybridge/backend/reverse"
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type UserInfo struct {
	Uname     string `json:"uname"`
	Uemail    string `json:"uemail"`
	Uphone    string `json:"uphone"`
	Upassword string `json:"upassword"`
}

func (a *App) Generate() UserInfo {
	rand.Seed(time.Now().UnixNano())

	nameLength := 5
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	name := make([]byte, nameLength)
	for i := range name {
		name[i] = letters[rand.Intn(len(letters))]
	}

	emailLength := 8
	localPart := make([]byte, emailLength)
	for i := range localPart {
		localPart[i] = letters[rand.Intn(len(letters))]
	}
	email := string(localPart) + "@gmail.com"

	phone := ""
	for i := 0; i < 11; i++ {
		phone += fmt.Sprintf("%d", rand.Intn(10))
	}
	config.Email = email
	config.Password = string(localPart)

	return UserInfo{
		Uname:     string(name),
		Uemail:    email,
		Uphone:    phone,
		Upassword: string(localPart),
	}
}
func (a *App) DeleteConfigFile() {
	configPath := os.Getenv("USERPROFILE") + "\\.cpolar\\cpolar1.yml"
	err := os.Remove(configPath)
	if err != nil {
		fmt.Println("删除文件失败:", err)
		return
	}
	fmt.Println("文件已成功删除:", configPath)
}

func (a *App) Start() {
	reverse.StartNc()
}
func (a *App) GetId() string {
	return reverse.GetId()
}

// 控制隧道的函数
var user *cpolar.User

func init() {
	//user = cpolar.Run()
}
func (a *App) GenerateTunnel() *cpolar.User {
	return cpolar.Run()
}
func (a *App) DeleteTunnel(id string) {
	user.Delete(id)
	user.Watch()
}
func (a *App) AddTunnel(port string, proto string) {
	user.Add(port, proto)
	user.Watch()
}

func (a *App) Stop() {

}
