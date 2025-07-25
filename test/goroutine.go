package main

import (
	"fmt"
	"time"
)

func main() {
	//写一个生产者消费者模型
	//定义一个通道
	ch := make(chan int) //无缓冲  ：发送和接收都会堵塞
	close(ch)
	//有缓冲
	//ch := make(chan int,10)
	//ch := make(chan string, 10)
	//启动一个协程充当生产者
	//go func(ch chan int) { //这里也是一个函数  ch是参数名  chan int 是参数类型
	//	i := 0
	//	for {
	//		//ch <- i 是 Go 语言中 向通道发送数据 的语法，属于通道（Channel）的核心操作之一。、
	//		ch <- i
	//		i++
	//		time.Sleep(time.Second)
	//		if i >= 10 {
	//			close(ch)
	//			break
	//		}
	//	}
	//}(ch) //传入ch参数
	//启动一个消费者
	go func(ch chan int) {
		//循环接收  循环结束时 表示通道已关闭且缓冲区为空
		for data := range ch {
			fmt.Println(data)
		}
		data := <-ch
		fmt.Println(data)
		fmt.Println("通道已关闭，退出接收循环")
	}(ch)

	time.Sleep(6 * time.Second)
	//select {}  一直阻塞

}
