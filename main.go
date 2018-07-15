package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Reader interface {
	Read(chan []byte)
}

type Writer interface {
	Write(chan string)
}

type LogProcess struct {
	rc    chan []byte
	wc    chan string
	read  Reader
	write Writer
}

type ReadFromFile struct {
	path string
}

type WriteToInfluxDB struct {
	influxDBDsn string
}

func (r *ReadFromFile) Read(rc chan []byte) {

	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file error:%s", err.Error()))
	}
	f.Seek(0, 2)
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(time.Second)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes error:%s", err.Error()))
		}
		rc <- line
	}
}

func (w *WriteToInfluxDB) Write(wc chan string) {
	for v := range wc {
		fmt.Print(v)
	}
}

func (l *LogProcess) Process() {
	// 127.0.0.1 - - [01/May/2018:12:29:56 +0800] "GET /img.zip HTTP/1.1" 200 3653244 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
	for v := range l.rc {
		l.wc <- strings.ToUpper(string(v))
	}
}

func main() {

	r := &ReadFromFile{
		path: "./access.log",
	}
	w := &WriteToInfluxDB{
		influxDBDsn: "username&password",
	}
	lp := &LogProcess{
		rc:    make(chan []byte),
		wc:    make(chan string),
		read:  r,
		write: w,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(time.Second * 30)
}
