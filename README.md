# 测试网速用

## 编译

```shell
go build main.go
```

## 实现的是以下两个shell脚本的功能

```shell
#!/bin/bash

# 监听指定端口 40001 用来发送数据
port=40001
echo "Sender is listening on port $port to send data..."

# 使用一个无限循环来保持服务端持续运行
while true; do
    # 每次连接后重新监听端口并发送数据
    echo "Waiting for a client to connect on port $port..."
    
    # 使用 nc 监听端口并发送 1MB 数据，发送完后关闭连接
    dd if=/dev/zero bs=1M count=1 | nc -l -p $port -q 1
    
    echo "1MB data sent to client. Waiting for next connection..."
done
```

```shell
#!/bin/bash

# 监听指定端口 40000 用来接收数据
port=40000
echo "Receiver is listening on port $port for incoming data..."

# 使用 nc 监听端口接收数据
while true; do
    # 监听端口并接收客户端的数据
    nc -l -p $port -q 1 | while read line; do
        echo "Received data: $line"
    done
done
```

## 测试是否成功

```shell
#!/bin/bash

# 服务端的 IP 地址
server_ip="127.0.0.1"
receive_port=40000  # 用于接收数据的服务端端口
send_port=40001     # 用于发送数据的服务端端口

# 向接收数据的服务端发送数据
echo "Sending data to receiver on port $receive_port"
exec 3<>/dev/tcp/$server_ip/$receive_port
echo "Hello, receiver!" >&3
exec 3>&-

# 计算接收数据的时间
start_time=$(date +%s%N)

# 向发送数据的服务端请求数据
echo "Requesting data from sender on port $send_port"
exec 3<>/dev/tcp/$server_ip/$send_port

# 从发送服务端接收数据，并将其输出到 /dev/null
dd bs=1M count=1 <&3 > /dev/null

end_time=$(date +%s%N)

# 计算下载时间
download_time=$((($end_time - $start_time)/1000000))  # 转换为毫秒
echo "Downloaded 1MB in $download_time ms"

# 关闭连接
exec 3>&-
```