/*
 * @Author: 肖和锋
 * @Date: 2022/6/29 15:22
 */
package work

import (
	"fmt"
	"github.com/buddyxiao/serial/ser"
	"testing"
)

func TestSaveToFile(t *testing.T) {
	file := "C:\\Users\\14640\\Desktop\\output.txt"
	mySerial := ser.NewMySerial()
	defer func() {
		mySerial.Close()
		fmt.Println("程序退出...")
	}()
	go mySerial.Run()
	saveToFile(file, mySerial)
}
