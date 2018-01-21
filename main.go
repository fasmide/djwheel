package main

import (
	"bytes"
	"time"

	"github.com/fasmide/djwheel/device"
	"github.com/fasmide/djwheel/plugins"

	_ "github.com/fasmide/djwheel/plugins/CPU"
	_ "github.com/fasmide/djwheel/plugins/Spectrum"
	_ "github.com/fasmide/djwheel/plugins/Volume"
)

func main() {

	// Initialize our physical usb device
	volumeEvents := make(chan int)
	device := device.NewDevice("/dev/ttyACM0", volumeEvents)
	go device.Loop()
	go plugins.HandleWheelEvents(volumeEvents)

	var buffer bytes.Buffer
	for {
		plugins.WriteTo(&buffer)
		device.Write(buffer.Bytes())
		buffer.Reset()
		time.Sleep(16 * time.Millisecond)

	}

}
