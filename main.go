// Program for captures live attendance events from the device, and logs those events

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chenall/gozk"
)

func main() {
	zkSocket := gozk.NewZK(0, "192.168.0.202", 4370, 0, "Asia/Saigon")
	if err := zkSocket.Connect(); err != nil {
		panic(err)
	}

	c := make(chan *gozk.Attendance, 50)

	// This function captures live attendance events and sends them to the channel.
	if err := zkSocket.LiveCapture(func(event *gozk.Attendance) {
		c <- event
	}); err != nil {
		panic(err)
	}

	//  logs the attendance events as they arrive.
	go func() {
		for event := range c {
			log.Println("event ", event)
		}
	}()

	gracefulQuit(zkSocket.StopCapture)
}

func gracefulQuit(f func()) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan

		log.Println("Stopping...")
		f()

		time.Sleep(time.Second * 1)
		os.Exit(1)
	}()

	for {
		time.Sleep(10 * time.Second) // or runtime.Gosched() or similar per @misterbee
	}
}
