package Spectrum

import (
	"io"
	"log"

	"github.com/fasmide/djwheel/audio"
	"github.com/fasmide/djwheel/plugins"
	colorful "github.com/lucasb-eyer/go-colorful"
)

type Spectrum struct {
	lastEvent *audio.SpectrumEvent
}

func init() {
	plugins.RegisterPlugin("spectrum", NewSpectrum())
}

func NewSpectrum() *Spectrum {
	s := &Spectrum{}
	// Initialize our audio input TODO: figure out this alsa output at runtime
	audioInput, err := audio.NewInput("alsa_output.pci-0000_00_1b.0.analog-stereo.monitor", 44100/60)

	if err != nil {
		log.Fatalf("Unable to open audio input: %s", err)
	}

	spectrum := audio.NewSpectrum(audioInput, 44100, false)
	spectrumUpdates := make(chan audio.SpectrumEvent)
	go spectrum.Loop(spectrumUpdates)

	go s.Loop(spectrumUpdates)
	return s
}

func (s *Spectrum) Priority() int {
	return 10
}

func (s *Spectrum) Write(to io.Writer) {

	if s.lastEvent == nil {
		return
	}

	for index := 0; index < 13; index++ {
		value := Map(s.lastEvent.Left.Power[index], 0, 90000, 0, 360)
		intens := Map(s.lastEvent.Left.Power[index], 0, 90000, 0, 1) * 2
		c := colorful.Hsv(value, 1, intens)
		to.Write([]byte{byte(c.R * 255), byte(c.G * 255), byte(c.B * 255)})
	}
	for index := 12; index >= 0; index-- {
		value := Map(s.lastEvent.Right.Power[index], 0, 90000, 0, 360)
		intens := Map(s.lastEvent.Right.Power[index], 0, 90000, 0, 1) * 2
		c := colorful.Hsv(value, 1, intens)
		to.Write([]byte{byte(c.R * 255), byte(c.G * 255), byte(c.B * 255)})
	}
}

// Loop just reads all spectrum events and saves the latest
func (s *Spectrum) Loop(updates chan audio.SpectrumEvent) {
	for e := range updates {
		s.lastEvent = &e

	}
}

func Map(x, inMin, inMax, outMin, outMax float64) float64 {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}