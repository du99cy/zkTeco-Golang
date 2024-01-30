package main
```
import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	gozk "github.com/0x19/gozk"
)
```

```
func main() {
	// Initialize ZKTeco device
	zkSocket := gozk.NewZK(0, "192.168.0.202", 4370, 0, "Asia/Saigon")

	// Connect to the ZKTeco device
	if err := zkSocket.Connect(); err != nil {
		panic(err)
	}

	// Create a buffered channel for attendance events
	c := make(chan *gozk.Attendance, 50)

	// Capture live attendance events and send them to the channel
	if err := zkSocket.LiveCapture(func(event *gozk.Attendance) {
		c <- event
	}); err != nil {
		panic(err)
	}

	// Log the attendance events as they arrive
	go func() {
		for event := range c {
			log.Println("event ", event)
		}
	}()

	// Handle graceful shutdown
	gracefulQuit(zkSocket.StopCapture)
}
```

```
// Graceful shutdown function
func gracefulQuit(f func()) {
	// Create a channel for receiving signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Goroutine for handling shutdown
	go func() {
		<-sigChan

		log.Println("Stopping...")
		f()

		time.Sleep(time.Second * 1)
		os.Exit(1)
	}()

	// Infinite loop to keep the program running
	for {
		time.Sleep(10 * time.Second)
	}
}
```
