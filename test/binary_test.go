package test

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Demo struct {
	N16 uint16
	N32 uint32
}

func ExampleBianry() {
	// 对象
	ss := &Demo{N16: 1, N32: 2}
	var bufferB, bufferL bytes.Buffer
	// 这里的binary的写入模式有2种:BigEndian 大端 , LittleEndian 小端; 假设都是写入1到内存,大端就是 01 而小端是 10
	binary.Write(&bufferB, binary.BigEndian, ss)
	// buffer大端结果: []uint8 len: 6, cap: 64, [0,1,0,0,0,2]
	binary.Write(&bufferL, binary.LittleEndian, ss)
	// buffer小端结果: []uint8 len: 6, cap: 64, [1,0,2,0,0,0]
	fmt.Println("\n大端bufferB: ", bufferB.Bytes())
	fmt.Println("小端bufferL: ", bufferL.Bytes())

	var un uint64 = 2
	var buffer_b bytes.Buffer
	var buffer_l bytes.Buffer
	binary.Write(&buffer_b, binary.BigEndian, un)    // []uint8 len: 8, cap: 64, [0,0,0,0,0,0,0,2]
	binary.Write(&buffer_l, binary.LittleEndian, un) // []uint8 len: 8, cap: 64, [2,0,0,0,0,0,0,0]

	fmt.Printf("n=%d大端: %v \n", un, buffer_b.Bytes()) // n=2大端: [0 0 0 0 0 0 0 2]
	fmt.Printf("n=%d小端: %v \n", un, buffer_l.Bytes()) // n=2小端: [2 0 0 0 0 0 0 0]

	var n1 int = 1
	var buffer_n1 bytes.Buffer
	binary.Write(&buffer_n1, binary.BigEndian, n1) //这个无效,因为int不能用于binary.Write
	fmt.Println("buffer_n1: ", buffer_n1)          //buffer_n1:  {[] 0 0}

	var n64 int64 = 1
	var buffer_n64 bytes.Buffer
	binary.Write(&buffer_n64, binary.BigEndian, n64)
	fmt.Println("buffer_n64: ", buffer_n64) //buffer_n64:  {[0 0 0 0 0 0 0 1] 0 0}

	var b1 bool = true
	var buffer_b1 bytes.Buffer
	binary.Write(&buffer_b1, binary.BigEndian, b1)
	fmt.Println("buffer_b1: ", buffer_b1) //buffer_b1:  {[1] 0 0}
	// output:
	// 大端bufferB:  [0 1 0 0 0 2]
	// 小端bufferL:  [1 0 2 0 0 0]
	// n=2大端: [0 0 0 0 0 0 0 2]
	// n=2小端: [2 0 0 0 0 0 0 0]
	// buffer_n1:  {[] 0 0}
	// buffer_n64:  {[0 0 0 0 0 0 0 1] 0 0}
	// buffer_b1:  {[1] 0 0}
}
