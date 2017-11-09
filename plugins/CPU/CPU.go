package cpu

import (
	"io"
	"log"
	"math"
	"time"

	"github.com/fasmide/djwheel/plugins"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/shirou/gopsutil/cpu"
)

type CPU struct {
	Cores []Core
}

type Core struct {
	color colorful.Color
	busy  float64
}

func init() {
	plugins.RegisterPlugin("cpu", NewCPU())
}

func NewCPU() *CPU {
	data, err := cpu.Times(true)

	if err != nil {
		panic("[CPU] unable to initialize cpu plugin")
	}

	c := &CPU{Cores: make([]Core, len(data), len(data))}

	for index, cpu := range data {
		calc := (360 / len(data)) * index
		color := colorful.Hsv(float64(calc), 1, 1)
		c.Cores[index].color = color
		c.Cores[index].busy = getAllBusy(cpu)

	}

	go c.Collect()
	return c
}

func (c *CPU) Collect() {
	var data []cpu.TimesStat
	var err error
	for {
		data, err = cpu.Times(true)
		if err != nil {
			log.Printf("[CPU] Unable to cputid")
			return
		}
		for i, cpu := range data {
			c.Cores[i].busy = getAllBusy(cpu) / 4
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (c *CPU) Priority() int {
	return 11
}

func (c *CPU) Write(to io.Writer) {

	b := make([]byte, 26*3)
	for _, core := range c.Cores {

		secPos := math.Mod(core.busy*100, 100)
		pos := int(math.Mod(secPos, 26)) * 3
		b[pos] = byte(core.color.R * 255)
		b[pos+1] = byte(core.color.G * 255)
		b[pos+2] = byte(core.color.B * 255)
	}
	to.Write(b)
}

func getAllBusy(t cpu.TimesStat) float64 {
	busy := t.User + t.System + t.Nice + t.Iowait + t.Irq +
		t.Softirq + t.Steal + t.Guest + t.GuestNice + t.Stolen
	return busy
}
