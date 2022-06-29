/*
 * @Author: 肖和锋
 * @Date: 2022/6/29 15:16
 */
package work

import (
	"fmt"
	"github.com/buddyxiao/serial/ser"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func saveToFile(saveFile string, ser *ser.MySerial) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	file, err := os.OpenFile(saveFile, os.O_CREATE|os.O_RDWR, 666)
	defer file.Close()
	if err != nil {
		log.Fatal("文件打开失败")
		return
	}
	data := ser.GetData()
	file.Write([]byte("lng,lat\n"))
	for true {
		output := fmt.Sprintf("%s,%s\n", data["lng"], data["lat"])
		fmt.Println("成功写入：", output)
		file.Write([]byte(output))
		time.Sleep(1 * time.Second)
		select {
		case <-c:
			goto Loop
		default:
		}
	}
Loop:
	fmt.Println("文件写入完毕")
}
