// 模拟TCP客户端
package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("TCP Client starts...")

	time.Sleep(1*time.Second)

	// 1. 连接TCP服务器，得到conn
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Printf("Client start: %s\n", err)
		return
	}

	for {
		// 2. 调用conn，向服务器Write数据
		if _, err = conn.Write([]byte("Hello, Zinx v0.1")); err != nil {
			fmt.Printf("Client: Write conn: %s\n", err)
			return
		}

		// 3. 接收服务器回传数据并打印
		buf := make([]byte, 512)
		if cnt, err := conn.Read(buf); err != nil {
			fmt.Printf("Client: Read conn: %s\n", err)
			return
		} else {
			fmt.Printf("Server call back: %s, cnt = %d\n", string(buf), cnt)
		}

		// 4. 睡眠
		time.Sleep(1*time.Second)
	}

}
