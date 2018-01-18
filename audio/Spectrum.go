package audio

import (
	"fmt"
	"math"
	"sync"

	"github.com/mjibson/go-dsp/spectral"
	"github.com/mjibson/go-dsp/window"
)

type Spectrum struct {
	input         *Input
	rate          float64
	logScale      bool
	left          []float64
	right         []float64
	pwelchOptions spectral.PwelchOptions
}

type SpectrumEvent struct {
	Left  SpectrumChannel
	Right SpectrumChannel
}

func (s *SpectrumEvent) HasData() bool {

	data := false
	wg := sync.WaitGroup{}

	check := func(c *SpectrumChannel) {
		for _, v := range c.Power {
			if v > 0 {
				data = true
				wg.Done()
				return
			}
		}
		wg.Done()
	}

	wg.Add(1)
	check(&s.Left)

	wg.Add(1)
	check(&s.Right)

	wg.Wait()

	return data

}

type SpectrumChannel struct {
	Power []float64
	Freqs []float64
}

func NewSpectrum(i *Input, rate int, logScale bool) *Spectrum {

	return &Spectrum{
		i,
		float64(rate),
		logScale,
		make([]float64, rate/60),
		make([]float64, rate/60),
		spectral.PwelchOptions{
			NFFT:      256,
			Window:    window.Hann,
			Scale_off: false,
		},
	}
}

func (s *Spectrum) Loop(eventChan chan SpectrumEvent) {
	var wg sync.WaitGroup
	var left, right SpectrumChannel
	for {

		err := s.input.Read(s.left, s.right)
		if err != nil {
			panic(err)
		}

		wg.Add(1)
		go func() {
			left.Power, left.Freqs = s.WriteSpectrum(s.left)
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			right.Power, right.Freqs = s.WriteSpectrum(s.right)
			wg.Done()
		}()

		wg.Wait()

		eventChan <- SpectrumEvent{Left: left, Right: right}
		//Render(left.Power, left.Freqs)
	}
}

func (s *Spectrum) WriteSpectrum(data []float64) ([]float64, []float64) {

	power, freqs := spectral.Pwelch(data, s.rate, &s.pwelchOptions)
	if s.logScale {
		for i, x := range power {
			if x < 1 {
				power[i] = 0
			} else {
				power[i] = 10 * math.Log10(x)
			}
		}
	}
	return power, freqs
}

func Render(powers, freqs []float64) {
	fmt.Print("\033c")
	for i, freq := range freqs {
		fmt.Printf("%4.0f: %6.3f ", freq, powers[i])
		for index := 0; int(powers[i]) > index; index++ {
			fmt.Print("☗")
		}
		fmt.Print("\n")
		if i > 13 {
			break
		}
	}
}
