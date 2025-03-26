package main

import (
	"Ybridge/backend/reverse"
	"context"
	"fmt"
	"math/rand"
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

	return UserInfo{
		Uname:     string(name),
		Uemail:    email,
		Uphone:    phone,
		Upassword: string(localPart),
	}
}

func (a *App) Start() {
	reverse.StartNc()
}
func (a *App) GetId() string {
	return reverse.GetId()
}
