package cpu

import (
	"io"
	"log"
	"math"
	"sync"
	"time"

	"github.com/fasmide/djwheel/plugins"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/shirou/gopsutil/cpu"
)

// CPU represents our current cpu color strip
type CPU struct {
	sync.RWMutex
	Strip     []colorful.Color
	CPUColors []colorful.Color
	FadeColor colorful.Color
}

func init() {
	plugins.RegisterPlugin("cpu", NewCPU())
}

// NewCPU initializes CPU plugin
func NewCPU() *CPU {
	data, err := cpu.Times(true)

	if err != nil {
		panic("[CPU] unable to initialize cpu plugin")
	}

	c := &CPU{
		Strip:     make([]colorful.Color, 26, 26),
		CPUColors: make([]colorful.Color, len(data), len(data)),
		FadeColor: colorful.LinearRgb(0, 0, 0),
	}

	for index := range data {

		rotation := (360 / len(data)) * index
		color := colorful.Hsv(float64(rotation), 1, 0.8)
		c.CPUColors[index] = color
	}

	go c.Collect()

	return c
}

// Collect collects cpu stats and renders colors
func (c *CPU) Collect() {
	var data []cpu.TimesStat
	var err error

	// These positions is used to determinane if we need to blend colors togeather
	// it is properly not a good way to do it as we cannot know if these positions
	// are from the current loop, or the previous
	lastPositions := make([]int, len(c.CPUColors), len(c.CPUColors))

	for {
		c.Lock()
		data, err = cpu.Times(true)
		if err != nil {
			log.Printf("[CPU] Unable to collect cpu usage")
			return
		}

		// By fading all colors agent black, we get a trail
		c.FadeAllToBlack()

		for i, cpu := range data {
			secPos := math.Mod((getAllBusy(cpu)/4)*100, 100)
			pos := int(math.Mod(secPos, 26))

			// If another cpu currently resides in this location, blend them togeather
			if intInSlice(pos, lastPositions) {
				c.Strip[pos] = c.Strip[pos].BlendHsv(c.CPUColors[i], 0.5)
			} else {
				c.Strip[pos] = c.CPUColors[i]
			}

			lastPositions[i] = pos
		}
		c.Unlock()
		time.Sleep(16 * time.Millisecond)
	}
}

// Priority is just fixed to 5
func (c *CPU) Priority() int {
	return 5
}

// FadeAllToBlack fades all colors a tiny bit against black
// leaving a tails for all color "dots"
func (c *CPU) FadeAllToBlack() {
	for i, color := range c.Strip {
		c.Strip[i] = color.BlendRgb(c.FadeColor, 0.15)
	}
}
func (c *CPU) WriteTo(to io.Writer) {
	c.RLock()
	for _, color := range c.Strip {

		to.Write([]byte{
			byte(color.R * 255),
			byte(color.G * 255),
			byte(color.B * 255),
		})
	}
	c.RUnlock()
}

// WheelEvent dosent do anything - we dont have a use for these
func (c *CPU) WheelEvent(_ int) {

}

// util to add all busy cpu time
func getAllBusy(t cpu.TimesStat) float64 {
	busy := t.User + t.System + t.Nice + t.Iowait + t.Irq +
		t.Softirq + t.Steal + t.Guest + t.GuestNice + t.Stolen
	return busy
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
