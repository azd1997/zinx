package net

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)


func TestDataPack(t *testing.T) {
	go mockServer()
	// 这里如果不睡，会很快执行到go mockClient!!，那么会由于server还没来得及监听，client由于找不到server就退出了
	time.Sleep(2*time.Second)
	go mockClient()

	select {}	// 阻塞，避免退出
}

func mockServer() {
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("listen error: ", err)
		return
	}



	for {
		fmt.Println("Server is waiting...")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}

		fmt.Println("server is ready to unpack")

		go testUnpack(conn)
	}
}

func testUnpack(conn net.Conn) (err error) {
	// 处理客户端请求
	// 第一次拆包，把head读出来
	// 第二次读，把data读出来
	// tcp是字节流的协议，前面接收一部分，后面再接收一部分，后面是从前面接收的截止处接着收的

	dp := NewDataPack()
	for {
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read error: ", err)
			return err
		}

		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack head error: ", err)
			return err
		}

		// 第二次读，根据head中的dataLen读取data内容
		if msgHead.GetDataLen() > 0 {
			// msg有实际数据
			msg := msgHead.(*Message)
			msg.Data = make([]byte, msg.GetDataLen())
			_, err = io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("unpack data error: ", err)
				return err
			}

			// 消息读取完毕，打印看看
			fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}
	}
}


func mockClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("client dial error: ", err)
		return
	}

	// 创建封包对象
	// 模拟TCP粘包过程：构建两个Message，合并起来一起打包发出
	//创建一个封包对象 dp
	dp := NewDataPack()

	//封装一个msg1包
	msg1 := &Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	msg2 := &Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client temp msg2 err:", err)
		return
	}

	//将sendData1，和 sendData2 拼接一起，组成粘包
	sendData1 = append(sendData1, sendData2...)

	//向服务器端写数据
	conn.Write(sendData1)
}