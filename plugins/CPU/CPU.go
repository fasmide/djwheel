package cpu

import (
	"fmt"
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
	pos   float64
	color colorful.Color
	idle  float64
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
		log.Printf("%d", calc)
		color := colorful.Hsv(float64(calc), 1, 1)
		c.Cores[index].color = color
		c.Cores[index].pos = float64(calc)
		c.Cores[index].idle = cpu.Idle

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
			diff := cpu.Idle - c.Cores[i].idle
			fmt.Printf("%d: %f\n", i, diff)
			c.Cores[i].pos += (diff)
			c.Cores[i].pos = math.Mod(c.Cores[i].pos, 26)
			c.Cores[i].idle = cpu.Idle
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

		pos := int(core.pos) * 3
		b[pos] = byte(core.color.R * 255)
		b[pos+1] = byte(core.color.G * 255)
		b[pos+2] = byte(core.color.B * 255)
	}
	to.Write(b)
}
