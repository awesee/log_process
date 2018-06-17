package main

import (
	"fmt"
	"strings"
	"time"
)

type Reader interface {
	Read(chan string)
}

type Writer interface {
	Write(chan string)
}
type LogProcess struct {
	rc    chan string
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

func (r *ReadFromFile) Read(rc chan string) {
	line := "message"
	rc <- line
}

func (w *WriteToInfluxDB) Write(wc chan string) {
	fmt.Println(<-wc)
}

func (l *LogProcess) Process() {
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
}

func main() {

	r := &ReadFromFile{
		path: "/tem/access.log",
	}
	w := &WriteToInfluxDB{
		influxDBDsn: "username&password",
	}
	lp := &LogProcess{
		rc:    make(chan string),
		wc:    make(chan string),
		read:  r,
		write: w,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(time.Second)
}
