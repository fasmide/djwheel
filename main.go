package main

import (
	"log"

	"github.com/fasmide/djwheel/audio"
	"github.com/fasmide/djwheel/device"

	"github.com/lucasb-eyer/go-colorful"
)

func main() {
	log.Printf("Hello")

	// Initialize our physical usb device
	volumeEvents := make(chan int)
	device := device.NewDevice("/dev/ttyACM0", volumeEvents)
	go device.Loop()
	go HandleVolume(volumeEvents)

	// Initialize our audio input TODO: figure out this alsa output at runtime
	audioInput, err := audio.NewInput("alsa_output.pci-0000_00_1b.0.analog-stereo.monitor", 44100/60)

	if err != nil {
		log.Fatalf("Unable to open audio input: %s", err)
	}

	spectrum := audio.NewSpectrum(audioInput, 44100, false)
	spectrumUpdates := make(chan audio.SpectrumEvent)
	go spectrum.Loop(spectrumUpdates)

	HandleSpectrum(device, spectrumUpdates)

}

func HandleVolume(e chan int) {
	for e := range e {
		log.Printf("Its time to turn this shit up to %d!", e)
	}
}

func HandleSpectrum(d *device.Device, updates chan audio.SpectrumEvent) {

	// Here be random selected values, tweeked while
	// listning to music and looking at the wheel
	for specUpdate := range updates {
		for index := 0; index < 13; index++ {
			value := Map(specUpdate.Left.Power[index], 0, 90000, 0, 360)
			intens := Map(specUpdate.Left.Power[index], 0, 90000, 0, 1) * 2
			c := colorful.Hsv(value, 1, intens)
			d.Write([]byte{byte(c.R * 255), byte(c.G * 255), byte(c.B * 255)})
		}
		for index := 12; index >= 0; index-- {
			value := Map(specUpdate.Right.Power[index], 0, 90000, 0, 360)
			intens := Map(specUpdate.Right.Power[index], 0, 90000, 0, 1) * 2
			c := colorful.Hsv(value, 1, intens)
			d.Write([]byte{byte(c.R * 255), byte(c.G * 255), byte(c.B * 255)})
		}
	}
}

func Map(x, inMin, inMax, outMin, outMax float64) float64 {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}
