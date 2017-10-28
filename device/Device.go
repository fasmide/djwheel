package device

import (
	"bufio"
	"io"
	"log"
	"strconv"

	"github.com/jacobsa/go-serial/serial"
)

// Device represents our djwheel
type Device struct {
	port      io.ReadWriteCloser
	eventChan chan int
}

// NewDevice returns a device api and connects to the serial port
func NewDevice(devicePath string, e chan int) *Device {
	// Set up options.
	options := serial.OpenOptions{
		PortName: devicePath,

		// BaudRate does nothing, when using emulated USB serial we should
		// be able to reach 12Mbit/sec
		BaudRate:        300,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("[device] serial.Open: %v", err)
	}

	return &Device{port: port, eventChan: e}
}

func (d *Device) Loop() {
	// scan the serial port and split on newline
	scanner := bufio.NewScanner(d.port)
	scanner.Split(bufio.ScanLines)

	var b []byte
	var i int
	var err error

	for scanner.Scan() {
		b = scanner.Bytes()
		i, err = strconv.Atoi(string(b))

		if err != nil {
			log.Printf("[Device.loop] I cannot possibly parse this robbish \"%s\": %s",
				string(b),
				err,
			)
		}

		d.eventChan <- i
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading device: %s", err)
	}
}

func (d *Device) Write(b []byte) (int, error) {
	return d.port.Write(b)
}
