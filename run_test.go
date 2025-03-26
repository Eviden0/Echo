package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestApp_Generate(t *testing.T) {
	rand.Seed(time.Now().UnixNano()) // 初始化随机数种子

	// 生成随机名称
	nameLength := 5
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	name := make([]byte, nameLength)
	for i := range name {
		name[i] = letters[rand.Intn(len(letters))]
	}

	// 生成随机电子邮件
	emailLength := 8
	localPart := make([]byte, emailLength)
	for i := range localPart {
		localPart[i] = letters[rand.Intn(len(letters))]
	}
	email := string(localPart) + "@gmail.com"

	// 生成随机电话号码
	phone := ""
	for i := 0; i < 11; i++ {
		phone += fmt.Sprintf("%d", rand.Intn(10)) // 生成随机数字
	}
	fmt.Println(string(name), phone, email, string(localPart))
}
