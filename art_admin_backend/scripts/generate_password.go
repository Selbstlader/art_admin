package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 默认密码
	password := "123456"

	// 如果提供了参数，使用参数作为密码
	if len(os.Args) > 1 {
		password = os.Args[1]
	}

	// 生成 bcrypt 哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("生成密码哈希失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("原密码: %s\n", password)
	fmt.Printf("哈希值: %s\n", string(hash))
	fmt.Printf("\n复制以下内容到 SQL 中:\n'%s'\n", string(hash))
}
