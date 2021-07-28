package main

import (
	"log"
	"time"

	"github.com/google/gousb"
)

const (
	slm25RequestType uint8  = gousb.ControlIn | gousb.ControlClass | gousb.ControlInterface
	slm25Request     uint8  = 0x01
	slm25Value       uint16 = 0x0105
)

// Reader USB data reader
type Reader struct {
	dev *gousb.Device
}

// NewReader creates a new USB data reader
func NewReader(dev *gousb.Device) *Reader {
	return &Reader{
		dev: dev,
	}
}

// Read reads the data from USB every 1 second, and writes the result into channel `out`
func (r *Reader) Read(stop <-chan struct{}) chan [2]float64 {
	out := make(chan [2]float64, 1000)

	go func() {
		defer close(out)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		response := make([]byte, 64)
		for true {
			select {
			case <-stop:
				return
			case t := <-ticker.C:
				_, err := r.dev.Control(slm25RequestType, slm25Request, slm25Value, 0, response)
				if err != nil {
					log.Printf("could not read data: %v\n", err)

					continue
				}

				out <- [2]float64{
					float64(t.Unix()),
					(float64(response[7])*256 + float64(response[8])) / 10,
				}
			}
		}
	}()

	return out
}
