package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/gousb"
)

const (
	slm25VendorID  gousb.ID = 0x10c4
	slm25ProductID gousb.ID = 0x82cd
)

func main() {
	ctx := gousb.NewContext()
	defer ctx.Close()

	dev, err := ctx.OpenDeviceWithVIDPID(slm25VendorID, slm25ProductID)
	if err != nil {
		log.Fatalf("could not open a device: %v\n", err)
	}
	defer dev.Close()

	dev.SetAutoDetach(true)

	_, interfaceDone, err := dev.DefaultInterface()
	if err != nil {
		log.Fatalf("%s.DefaultInterface(): %v\n", dev, err)
	}
	defer interfaceDone()

	apiKey := os.Getenv("DD_API_KEY")
	if apiKey == "" {
		log.Fatalf("'DD_API_KEY' is not set\n")
	}

	metricName := os.Getenv("DD_METRIC_NAME")
	if metricName == "" {
		log.Fatalf("'DD_METRIC_NAME' is not set\n")
	}

	stop := make(chan struct{})
	submitterDone := NewSubmitter(apiKey, metricName).Submit(
		NewReader(dev).Read(stop),
	)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("terminating...")
		close(stop)
	}()

	<-submitterDone
	log.Println("terminated")
}
