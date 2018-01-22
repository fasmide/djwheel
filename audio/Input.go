package audio

import (
	"encoding/binary"
	"fmt"
	"io"
	"os/exec"
)

type Input struct {
	fd     io.ReadCloser
	buffer []int16
}

func NewInput(device string, bufferSize int) (*Input, error) {

	//pacat --record -d alsa_output.pci-0000_00_1b.0.analog-stereo.monitor
	// defaults to  rate 44100, signed-integer, little-endian, 16-bit and stereo
	cmd := exec.Command("pacat", "--record", "--latency-msec", "16", "-d", device)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to get pipe from pacat: %s", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("Unable to start pacat: %s", err)
	}

	return &Input{fd: pipe, buffer: make([]int16, bufferSize*2)}, nil
}

func (i *Input) Read(left, right []float64) error {

	if err := binary.Read(i.fd, binary.LittleEndian, i.buffer); err != nil {
		return err
	}
	for index := 0; index < len(i.buffer)/2; index++ {
		left[index] = float64(i.buffer[index*2])
		right[index] = float64(i.buffer[index*2+1])
	}
	return nil
}
