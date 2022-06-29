/*
 * @Author: 肖和锋
 * @Date: 2022/6/29 10:59
 */
package main

import (
	"fmt"
	"github.com/buddyxiao/serial/ser"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	mySerial := ser.NewMySerial()
	defer func() {
		mySerial.Close()
		fmt.Println("程序退出...")
	}()
	go mySerial.Run()
	go readData(mySerial)
	select {
	case <-c:
	}
}

// 结果：map[lat:30.20238005 lng:112.12500994]
func readData(mySerial *ser.MySerial) {
	for true {
		data := mySerial.GetData()
		fmt.Println(data)
		time.Sleep(1 * time.Second)
	}
}
