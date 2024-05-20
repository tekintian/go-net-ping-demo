package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"time"
)

var (
	timeout = flag.Int64("w", 1000, "timeout for ping requests")
	size    = flag.Int("l", 32, "size of the packet buffer")
	count   = flag.Int("n", 4, "number of send requests")
	rtype   = flag.String("rt", "ip:icmp", "the type of request to send") // ip:icmp   tcp  udp
)

// uint8 取值范围: 0 -- 255.
// uint16 取值范围 0 -- 65535

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	ID          uint16
	SequenceNum uint16
}

func help() {
	// 打印多行使用帮助信息
	fmt.Print(`
用法: go-ping [-n count] [-l size] [-w timeout] [-rt 请求类型] target_name
选项:
    -l size        发送缓冲区大小。
    -n count       要发送的回显请求数。
    -w timeout     等待每次回复的超时时间(毫秒)。
    -rt ip:icmp   要发送的请求类型 默认 ip:icmp, 注意ip:icmp请求在mac里面有安全限制不允许第三方发送

`)
}
func main() {
	// os.Args
	flag.Parse()
	fmt.Printf("timeout:%v ,size:%v count: %v \n", *timeout, *size, *count)
	target := "www.baidu.com" // 默认目标
	if len(os.Args) > 1 {
		target = os.Args[len(os.Args)-1]
	} else {
		help() // 显示帮助信息
	}
	// net dial timeout
	conn, err := net.DialTimeout(*rtype, target, time.Duration(*timeout)*time.Millisecond)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	fmt.Printf("正在ping %s [%s] 具有 %d 字节的数据", target, conn.RemoteAddr().String(), *size)

	var SuccessTimes int64             // 成功次数
	var FailTimes int64                // 失败次数
	var minTime = int64(math.MaxInt64) // math.MaxInt64 返回和还是一个int类型的值,所以这里强转int6
	var maxTime int64
	var totalTime int64

	// 循环指定的次数
	for i := 0; i < *size; i++ {

		t1 := time.Now() //记录开始时间

		icmp := &ICMP{
			Type:        8, //请求头部信息, 8表示为icmp请求 ping请求
			Code:        0,
			Checksum:    0,
			ID:          1,
			SequenceNum: uint16(1),
		}

		data := make([]byte, *size)
		var buffer bytes.Buffer
		// 这里的binary的写入模式有2种:BigEndian 大端 , LittleEndian 小端; 假设都是写入1到内存,大端就是 01 而小端是 10
		binary.Write(&buffer, binary.BigEndian, icmp) // 写入icmp包头写入到buffer
		// buffer大端结果: []uint8 len: 8, cap: 64, [8,0,0,0,0,1,0,1]
		// buffer小端结果: []uint8 len: 8, cap: 64, [8,0,0,0,1,0,1,0]
		fmt.Println("\nbuffer: ", buffer.Bytes()) //output: []uint8 len: 8, cap: 64, [8,0,0,0,1,0,1,0]
		buffer.Write(data)                        // 写入32位空数据
		data = buffer.Bytes()                     // 重新赋值data
		fmt.Println("data: ", data)               //output: []uint8 len: 40, cap: 64, [8,0,0,0,0,1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]
		checkSum := checkSum(data) // output: 63485
		// 设置 checksum
		data[2] = byte(checkSum >> 8) // uint16 占2个字节 , 这里将checkSum 右移8位即取出他的高8位
		data[3] = byte(checkSum)

		// 设置sequenceNum
		data[6] = byte(icmp.SequenceNum >> 8)
		data[7] = byte(icmp.SequenceNum)

		conn.SetDeadline(time.Now().Add(time.Duration(*timeout) * time.Second))
		n, err := conn.Write(data)
		if err != nil {
			log.Fatalln(err, n)
			FailTimes++
			continue // 转到下一个操作
		}

		// 定义一个用于读取数据的buf
		buf := make([]byte, 65535) // 这里设置为0的话拿不到数据
		nr, err := conn.Read(buf)  // 从链接中读取数据 一次性读
		if err != nil {
			log.Fatalln(err, nr)
			FailTimes++
			continue
		}
		/* //循环读取数据
		buf := make([]byte, 0, 512)
		for {
			nr, err := conn.Read(buf[len(buf):cap(buf)]) // 从链接中读取数据
			buf = buf[:len(buf)+nr]
			if err != nil {
				break
			}
			fmt.Println(buf)
			if len(buf) == cap(buf) { //扩容
				buf = append(buf, 0)[:len(buf)]
			}
		}
		*/
		ms := time.Since(t1).Milliseconds()

		if minTime > ms {
			minTime = ms
		}
		if maxTime < ms {
			maxTime = ms
		}
		totalTime += ms
		SuccessTimes++
		// icmp响应包的第 12至15位为IP地址
		fmt.Printf("来自 %d.%d.%d.%d 的回复：字节=%d 时间=%dms TTL=%d \n", buf[12], buf[13], buf[14], buf[15], n-28, ms, buf[8])

		time.Sleep(1 * time.Second) // sleep for 1 second

	}

}

func checkSum(data []byte) uint16 {
	length := len(data) // 获取数据长度
	index := 0
	var sum uint32 = 0
	for length > 1 {
		// 俩俩拼接并求和
		sum += uint32(data[index])<<8 + uint32(data[index+1]) // 相连2个字节相加, data[index]<<8将前面的数左移8为 这样相加后就是16位
		length -= 2
		index += 2
	}

	if length != 0 {
		sum += uint32(data[index]) // 最后多余的数,直接链接到求和计算中
	}

	high16 := sum >> 16 // 取出sum的高16位  ( 将sum右移16位)
	if high16 != 0 {    // 如果高16位不为0, 那就重复的将高16位和低16位相加,直到高16位为0
		sum = high16 + uint32(uint16(data[index])) // 由于高位数转低位数,会拿到低位的值,超出的高位会被丢弃, 所以这里 uint16(data[index]) 获取到的就是sum的低位数, 然后在转换为 uint32
		high16 = sum >> 16                         // 再次 取出sum的高16位的值  如果结果为0 这验证成功
	}

	return uint16(^sum) // ^sum值 按位异或 (正数按位异或的结果为+1后的负数) 然后转换为 uint16
}

// 另外一个 CheckSum 计算校验和 算法, 功能和上面的是一样的,更精简一些而已
func CheckSumOther(data []byte) uint16 {
	var sum uint32
	var length = len(data)
	var index int

	for length > 1 { // 溢出部分直接去除
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length == 1 {
		sum += uint32(data[index])
	}
	// uint16(sum>>16) 取出sum的高位数, uint16(sum) 取出sum的低位数
	sum = uint32(uint16(sum>>16) + uint16(sum))
	// []uint8 40位的数据[8,0,0,0,0,1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]
	// sum: 2050 = 0x802  取反后^sum 为 -2051 = 0x0 转换为uint16(^sum) 结果 63485 = 0xf7fd
	return uint16(^sum)
}
