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
	webcam, _ := gocv.VideoCaptureDevice(0)
	img := gocv.NewMat()
	magenta := color.RGBA{255, 0, 255, 0}
	turquoise := color.RGBA{64, 224, 208, 0}

	for {
		webcam.Read(&img)
		tmpl := gocv.IMRead("images/gopher_template.png", gocv.IMReadColor)
		m := gocv.NewMat()
		result := gocv.NewMat()
		gocv.MatchTemplate(img, tmpl, &result, gocv.TmCcoeffNormed, m)
		m.Close()
		_, maxConfidence, _, maxLoc := gocv.MinMaxLoc(result)
		r := image.Rect(maxLoc.X, maxLoc.Y, maxLoc.X+tmpl.Cols(), maxLoc.Y+tmpl.Rows())
		if maxConfidence < 0.50 {
			fmt.Print("Move a bit... Confidence level too low!!!:", maxConfidence)
		} else {
			gocv.Rectangle(&img, r, turquoise, 3)
			size := gocv.GetTextSize("That's a Gopher", gocv.FontHersheyPlain, 4, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, "That's a Gopher", pt, gocv.FontHersheyPlain, 4, magenta, 2)
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
