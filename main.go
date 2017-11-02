package main

import (
	"bytes"
	"log"
	"time"

	"github.com/fasmide/djwheel/device"
	"github.com/fasmide/djwheel/plugins"

	_ "github.com/fasmide/djwheel/plugins/CPU"
	_ "github.com/fasmide/djwheel/plugins/Spectrum"
)

func main() {
	log.Printf("Hello")

	// Initialize our physical usb device
	volumeEvents := make(chan int)
	device := device.NewDevice("/dev/ttyACM0", volumeEvents)
	go device.Loop()
	go HandleVolume(volumeEvents)

	var buffer bytes.Buffer
	for {
		plugins.Write(&buffer)
		device.Write(buffer.Bytes())
		buffer.Reset()
		time.Sleep(16 * time.Millisecond)

	}

}

func HandleVolume(e chan int) {
	for e := range e {
		log.Printf("Its time to turn this shit up to %d!", e)
	}
}
