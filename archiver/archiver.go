package archiver

import (
	"time"
)

func NewArchiver(logPath string) *Archiver {
	return &Archiver{
		logPath:   logPath,
		workQueue: make(chan string),
	}
}

type Archiver struct {
	logPath   string
	workQueue chan string
}

func (a *Archiver) Boot() {
	go a.startScanner()
	go a.startConsumer()
}

func (a *Archiver) startConsumer() {
	for task := range a.workQueue {
		a.gzipFile(task)
	}
}

func (a *Archiver) startScanner() {
	for {
		a.scanLogPath()
		time.Sleep(time.Second * 60)
	}
}
