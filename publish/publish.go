package main

import (
	"bufio"
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
	nopts.MaxPayload = 2764800
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
	fmt.Println("Press ESC button or Ctrl-C to exit this program")
	fmt.Println("Press RETURN to send webcam image")

	for {
		consoleReader := bufio.NewReaderSize(os.Stdin, 1)
		fmt.Print(">> PRESS RETURN")
		input, _ := consoleReader.ReadByte()

		ascii := input

		if ascii == 27 || ascii == 3 {
			fmt.Println("Exiting...")
			os.Exit(0)
		}

		if ascii == 10 {
			if ok := webcam.Read(&img); !ok {
				fmt.Printf("cannot read device %v\n", 0)
				return
			}
			publish(sc, img.ToBytes())
			fmt.Println("---> IMAGE SENT")
		}
	}
}

func publish(sc stan.Conn, data []byte) {
	if err := sc.Publish("foo", data); err != nil {
		log.Fatalln(err)
	}
}
