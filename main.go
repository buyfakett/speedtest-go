package main

import (
    "fmt"
    "log"
    "net"
    "bufio"
)

var (
    sender_port   = "40000" // 发送端口
    receiver_port = "40001" // 接收端口
)

func main() {
    // 启动接收数据的服务器
    go startReceiver()

    // 启动发送数据的服务器
    go startSender()

    // 阻塞主程序直到按下 Ctrl+C
    select {}
}

// 接收数据并打印的功能
func startReceiver() {
    fmt.Printf("Receiver is listening on port %s for incoming data...\n", sender_port)

    // 监听指定端口
    listen, err := net.Listen("tcp", ":"+sender_port)
    if err != nil {
        log.Fatalf("Error starting receiver server: %v", err)
    }
    defer listen.Close()

    for {
        // 接受连接
        conn, err := listen.Accept()
        if err != nil {
            log.Printf("Error accepting connection: %v", err)
            continue
        }

        // 启动一个 goroutine 来处理每个连接
        go handleReceiverConnection(conn)
    }
}

func handleReceiverConnection(conn net.Conn) {
    defer conn.Close()
    scanner := bufio.NewScanner(conn)

    // 读取传输的数据
    for scanner.Scan() {
        fmt.Printf("Received data: %s\n", scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Printf("Error reading data: %v", err)
    }
}

// 发送 1MB 数据的功能
func startSender() {
    fmt.Printf("Sender is listening on port %s to send data...\n", receiver_port)

    // 监听指定端口
    listen, err := net.Listen("tcp", ":"+receiver_port)
    if err != nil {
        log.Fatalf("Error starting sender server: %v", err)
    }
    defer listen.Close()

    for {
        // 接受连接
        conn, err := listen.Accept()
        if err != nil {
            log.Printf("Error accepting connection: %v", err)
            continue
        }

        // 启动一个 goroutine 来发送数据
        go sendData(conn)
    }
}

func sendData(conn net.Conn) {
    defer conn.Close()

    // 生成 1MB 的数据（模拟发送）
    data := make([]byte, 1024*1024) // 1MB 的零字节数据

    // 发送数据
    _, err := conn.Write(data)
    if err != nil {
        log.Printf("Error sending data: %v", err)
        return
    }

    fmt.Println("1MB data sent to client. Waiting for next connection...")
}

