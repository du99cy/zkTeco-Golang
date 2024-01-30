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

func newMachine(host string) *gozk.ZK {
	return gozk.NewZK(
		0,
		host,
		4370,
		0,
		"Asia/Saigon")
}

func main() {
	machine := newMachine("192.168.1.204")
	if err := machine.Connect(); err != nil {
		panic(err)
	}

	c := make(chan *gozk.Attendance, 50)

	// This function captures live attendance events and sends them to the channel.
	if err := machine.LiveCapture(func(event *gozk.Attendance) {
		c <- event
	}); err != nil {
		panic(err)
	}

	//  logs the attendance events as they arrive.
	go func() {
		for event := range c {
			log.Println("ATTENDANCE LOG DATA: ", event)
		}
	}()

	gracefulQuit(machine.StopCapture)
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
