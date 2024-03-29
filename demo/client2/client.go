// 模拟TCP客户端
package main

import (
	"fmt"
	"io"
	"net"
	"time"

	enet "github.com/azd1997/zinx/net"
)

func main() {
	fmt.Println("TCP Client starts...")

	// 给服务端开启服务的时间
	time.Sleep(1 * time.Second)

	// 1. 连接TCP服务器，得到conn
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Printf("Client start: %s\n", err)
		return
	}

	for {

		// 封包消息
		dp := enet.NewDataPack()
		msg, _ := dp.Pack(enet.NewMessage(1,[]byte("Zinx V1.0 Client Test Message")))
		_, err := conn.Write(msg)
		if err !=nil {
			fmt.Println("write error err ", err)
			return
		}

		//先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}

		//将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*enet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		// 4. 睡眠
		time.Sleep(1 * time.Second)
	}

}
