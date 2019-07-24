package main

import (
	"fmt"
	_ "image/jpeg"
	"log"
	"os"
	"strings"

	stand "github.com/nats-io/nats-streaming-server/server"
	"github.com/nats-io/stan.go"
	"gocv.io/x/gocv"
)

func main() {
	opts := stand.GetDefaultOptions()
	nopts := stand.NewNATSOptions()
	nopts.MaxPayload = 30000000
	s, _ := stand.RunServerWithOpts(opts, nopts)
	defer s.Shutdown()
	sc, err := stan.Connect(
		"test-cluster",
		"client-1",
		stan.Pings(1, 3),
		stan.NatsURL(strings.Join(os.Args[1:], ",")),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer sc.Close()
	cv(sc)
}

func cv(sc stan.Conn) {
	webcam, _ := gocv.VideoCaptureDevice(0)
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		go publish(sc, img.ToBytes())
		fmt.Println("---> IMAGE SENT")
	}
}

func publish(sc stan.Conn, data []byte) {
	if err := sc.Publish("video", data); err != nil {
		log.Fatalln(err)
	}
}
