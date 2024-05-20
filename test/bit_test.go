package test

import "fmt"

func ExampleDemo() {
	// uint32 Range: 0 through 4294967295.
	// uint16 Range: 0 through 65535
	var sum uint32 = 65635
	high16 := uint16(sum >> 16) // 32位的sum右移16位 然后转换为uint16, 即获取sum的16位高位数
	low16 := uint16(sum)        // 直接将32位的数转换为16位数, 高位会被丢弃,即获得了 sum 的低16位数
	fmt.Printf("high16: %v, low16: %v sum>>16结果为:%v \n", high16, low16, sum >> 16)
	// 单个正数 按位异或  ^100 结果为加1后的负数 即 -101

	//output: high16: 1, low16: 99 sum>>16结果为:1
}
