/*
 * @Author: 肖和锋
 * @Date: 2022/6/29 11:10
 */
package ser

import (
	"bufio"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
	"strings"
	"sync"
)

type MySerial struct {
	Option       serial.OpenOptions
	Ser          io.ReadWriteCloser // 串口读写对象
	Temp         chan string        // 临时存放原始定位数据
	Data         map[string]string  // 保存解析后实时的定位信息
	FilterOption string             // 过滤数据的字段
	sync.RWMutex
}

func NewMySerial() *MySerial {
	return &MySerial{
		Option: serial.OpenOptions{
			PortName:        "COM5",
			BaudRate:        115200,
			DataBits:        8,
			StopBits:        1,
			MinimumReadSize: 4,
		},
		Ser: nil, Data: make(map[string]string),
		FilterOption: "$GPRMC",
		Temp:         make(chan string),
	}
}

// Open 打开串口
func (s *MySerial) Open() io.ReadWriteCloser {
	open, err := serial.Open(s.Option)
	if err != nil {
		log.Fatal("串口打开失败!")
	}
	s.Ser = open
	return open
}

// Close 关闭串口
func (s *MySerial) Close() {
	if s.Ser == nil {
		return
	}
	s.Ser.Close()
}

// Run 运行串口程序
func (s *MySerial) Run() {
	s.Open()
	go s.ReadFromSerial() // 读取数据
	go s.writerToData()   // 保存数据
}

// ReadFromSerial 不停地从串口中读取数据，并保存到Data字段中
func (s *MySerial) ReadFromSerial() {
	for true {
		reader := bufio.NewReader(s.Ser)
		readString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("reader.ReadString err: ", err)
			return
		}
		if strings.HasPrefix(readString, s.FilterOption) {
			s.Temp <- readString
		}
	}
}

// GetData 获取到经纬度信息     结果：map[lat:30.20238005 lng:112.12500994]
func (s *MySerial) GetData() map[string]string {
	for true {
		if _, ok := s.Data["lng"]; !ok {
			fmt.Println("数据还未到...")
		} else {
			break
		}
	}
	s.RLock()
	data := s.Data
	s.RUnlock()
	return data
}

// 保存经纬度信息
func (s *MySerial) writerToData() {
	for true {
		str := <-s.Temp
		if data, ok := s.processStr(str); ok {
			s.Lock()
			s.Data["lng"] = data[0]
			s.Data["lat"] = data[1]
			s.Unlock()
		}
	}

}

// 预处理数据，取出经纬度 $GPRMC,042510.00,A,3020.242890,N,11212.499549,E,0.0,329.9,290622,,,A*59
func (s *MySerial) processStr(readString string) ([]string, bool) {
	if strings.TrimSpace(readString) == "" {
		return nil, false
	}
	splitStr := strings.Split(readString, ",")
	var justify = splitStr[2]
	if justify == "A" {
		var data = make([]string, 2)
		var lng = movePoint(splitStr[5], 2) // 经度
		var lat = movePoint(splitStr[3], 2) // 维度
		data[0], data[1] = lng, lat
		return data, true
	} else {
		return nil, false
	}
}

// 3020.242890
func movePoint(raw string, num int) string {
	pIdx := strings.IndexRune(raw, '.')
	sArr := []byte(raw)
	for i := pIdx; i > pIdx-num; i-- {
		temp := sArr[i-1]
		sArr[i-1] = sArr[i]
		sArr[i] = temp
	}
	return string(sArr)
}
