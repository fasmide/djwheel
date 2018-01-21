package Volume

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/fasmide/djwheel/plugins"
	colorful "github.com/lucasb-eyer/go-colorful"
)

func init() {
	plugins.RegisterPlugin("volume", NewVolume())
}

const inputtimeout = time.Second * 1
const volumetimeout = time.Millisecond * 50

type Volume struct {
	sync.RWMutex
	// currentVolume: 0.0 == muted 1.0 == 100% volume
	currentVolume   float64
	currentPosition int
	rendering       bool
	// when this fires - we will stop rendering
	inputTimeout *time.Timer

	// when this fires - we will set the volume
	volumeTimeout *time.Timer
}

func NewVolume() *Volume {
	// TODO: figure out what the current volume is instead of hardcoding 0.25
	// pactl get-sink-volume alsa_output.pci-0000_00_1f.3.analog-stereo 0.05
	// hah! we dont need to know the sinks name, we can just use @DEFAULT_SINK@
	v := Volume{
		currentVolume: 0.25,
		volumeTimeout: time.NewTimer(volumetimeout),
		inputTimeout:  time.NewTimer(inputtimeout),
	}

	// start goroutine to handle timeouts
	go v.handleTimeouts()

	return &v

}

func (v *Volume) handleTimeouts() {
	for {
		select {
		case <-v.inputTimeout.C:
			v.inputTimeout.Stop()
			v.Lock()
			v.rendering = false
			v.Unlock()
		case <-v.volumeTimeout.C:
			v.setSystemVolume()
		}
	}
}

func (v *Volume) setSystemVolume() {
	v.RLock()
	cmd := exec.Command("pactl",
		"set-sink-volume",
		"@DEFAULT_SINK@",
		fmt.Sprintf("%.0f%%", v.currentVolume*100),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("unable to change volume: %s: %s", err, output)
	}
	v.RUnlock()
}

func (v *Volume) Priority() int {
	v.RLock()
	defer v.RUnlock()
	if v.rendering {
		return 11
	}
	return 0
}

func (v *Volume) WheelEvent(pos int) {

	v.Lock()

	// reset both timers
	// Note: when calling stop the timer could have already
	// been fired - but in our case it does not really matter
	v.inputTimeout.Stop()
	v.inputTimeout.Reset(inputtimeout)

	v.volumeTimeout.Stop()
	v.volumeTimeout.Reset(volumetimeout)

	// change the volume
	if pos < v.currentPosition && v.currentVolume >= 0 {
		v.currentVolume -= 0.002
	}
	if pos > v.currentPosition && v.currentVolume <= 1 {
		v.currentVolume += 0.002
	}

	// save the current position for the next event
	v.currentPosition = pos

	// enable rendering (well .. even if it already was...)
	v.rendering = true
	v.Unlock()
}

func (v *Volume) WriteTo(w io.Writer) {
	v.RLock()
	for i := 0; i < 26; i++ {
		if v.currentVolume*26 >= float64(i) {
			c := colorful.Hsv(v.currentVolume*150, 1, 0.5)
			w.Write([]byte{byte(c.R * 255), byte(c.G * 255), byte(c.B * 255)})
			continue
		}
		w.Write([]byte{0x00, 0x00, 0x00})

	}
	v.RUnlock()
}

// time, begining value, change in value, duration
func (v *Volume) easeOutCubic(t, b, c, d int) int {
	t = t/d - 1
	return c*(t*t*t+1) + b
}
