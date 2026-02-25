package main

import (
	"fmt"
	"net/http"
)

// main 函数是程序的入口，类似于 C/C++ 或 Java 的 main
// 在 Go 中，可执行程序必须属于 'main' 包
func main() {
	// 定义一个路由处理函数
	// 类似于 Express.js 的 app.get('/', (req, res) => { ... })
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// w (ResponseWriter) 用于写入响应数据
		// r (Request) 包含请求信息
		fmt.Fprintf(w, "Hello, go-zero learner!")
	})

	fmt.Println("Server starting on http://localhost:8080 ...")

	// 启动 HTTP 服务，监听 8080 端口
	// nil 表示使用默认的路由复用器 (DefaultServeMux)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
