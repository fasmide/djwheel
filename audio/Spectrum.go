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
	sync.RWMutex
	Left  SpectrumChannel
	Right SpectrumChannel
}

func (s *SpectrumEvent) HasData() bool {
	s.RLock()
	defer s.RUnlock()

	wg := sync.WaitGroup{}

	check := func(c *SpectrumChannel, data *bool) {
		for _, v := range c.Power {
			if v > 0 {
				*data = true
				wg.Done()
				return
			}
		}
		wg.Done()
	}

	left := false
	wg.Add(1)
	go check(&s.Left, &left)

	right := false
	wg.Add(1)
	go check(&s.Right, &right)

	wg.Wait()
	return left || right

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

func (s *Spectrum) Loop() *SpectrumEvent {
	var wg sync.WaitGroup
	var event SpectrumEvent

	go func() {
		for {
			err := s.input.Read(s.left, s.right)
			if err != nil {
				panic(err)
			}

			// event is shared with the spectrum plugin
			// we must ensure its not reading from this event
			event.Lock()
			wg.Add(1)
			go func() {
				event.Left.Power, event.Left.Freqs = s.WriteSpectrum(s.left)
				wg.Done()
			}()

			wg.Add(1)
			go func() {
				event.Right.Power, event.Right.Freqs = s.WriteSpectrum(s.right)
				wg.Done()
			}()

			wg.Wait()
			event.Unlock()
			//Render(left.Power, left.Freqs)
		}
	}()

	return &event
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
