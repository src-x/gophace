package main

import (
	"fmt"
	"image"
	"image/color"
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
	blue := color.RGBA{0, 0, 255, 0}

	cascadePath := "xml_files/haarcascade_frontalface_default.xml"
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(cascadePath) {
		fmt.Printf("Error reading cascade file")
		return
	}

	webcam, _ := gocv.VideoCaptureDevice(0)
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		for _, r := range rects {
			gocv.Rectangle(&img, r, blue, 3)

			size := gocv.GetTextSize("A Face", gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, "A Face", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
		}
		go publish(sc, img.ToBytes())
		fmt.Println("---> IMAGE SENT")
	}
}

func publish(sc stan.Conn, data []byte) {
	if err := sc.Publish("video", data); err != nil {
		log.Fatalln(err)
	}
}
