package Spectrum

import (
	"io"
	"log"

	"github.com/fasmide/djwheel/audio"
	"github.com/fasmide/djwheel/plugins"
	colorful "github.com/lucasb-eyer/go-colorful"
)

type Spectrum struct {
	spectrum *audio.SpectrumEvent
}

func init() {
	plugins.RegisterPlugin("spectrum", NewSpectrum())
}

func NewSpectrum() *Spectrum {
	s := &Spectrum{}
	// Initialize our audio input TODO: figure out this alsa output at runtime
	audioInput, err := audio.NewInput("alsa_output.pci-0000_00_1f.3.analog-stereo.monitor", 44100/60)

	if err != nil {
		log.Fatalf("Unable to open audio input: %s", err)
	}

	spectrum := audio.NewSpectrum(audioInput, 44100, false)
	s.spectrum = spectrum.Loop()

	return s
}

func (s *Spectrum) Priority() int {
	if s.spectrum == nil {
		return 1
	}

	// hasdata does its own read Lock
	// yeah it's gone mad this locking fix
	if s.spectrum.HasData() {
		return 10
	}

	return 1
}

func (s *Spectrum) WriteTo(to io.Writer) {
	s.spectrum.RLock()
	defer s.spectrum.RUnlock()
	// we must ensure we are not writing new spectrum data into this Spectrum
	// struct - woops this definitely got out of hand
	if s.spectrum == nil {
		return
	}

	for index := 0; index < 13; index++ {
		value := Map(float64(index), 0, 12, 0, 320)
		intens := Map(s.spectrum.Left.Power[index], 0, 90000, 0, 0.8) * 2
		c := colorful.Hsv(value, 1, intens)
		to.Write([]byte{byte(c.R * 255), byte(c.G * 255), byte(c.B * 255)})
	}
	for index := 12; index >= 0; index-- {
		value := Map(float64(index), 0, 12, 0, 320)
		intens := Map(s.spectrum.Right.Power[index], 0, 90000, 0, 0.8) * 2
		c := colorful.Hsv(value, 1, intens)
		to.Write([]byte{byte(c.R * 255), byte(c.G * 255), byte(c.B * 255)})
	}
}

// WheelEvent is just a placeholder for wheel events that
// we dont use for anything
func (s *Spectrum) WheelEvent(_ int) {

}

func Map(x, inMin, inMax, outMin, outMax float64) float64 {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}
