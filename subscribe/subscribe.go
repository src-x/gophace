package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten"
	stan "github.com/nats-io/stan.go"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

var camimg *ebiten.Image
var canvasImg *ebiten.Image

func init() {
	canvasImg, _ = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
	canvasImg.Fill(color.White)
}
func subscribe(c chan []byte, sc stan.Conn) {
	_, err := sc.Subscribe("video", func(m *stan.Msg) {
		fmt.Print(".")
		c <- m.Data
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	sc, err := stan.Connect(
		"test-cluster",
		"client-2",
		stan.Pings(1, 3),
		stan.NatsURL(strings.Join(os.Args[1:], ",")),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer sc.Close()
	imageChannel := make(chan []byte)
	go subscribe(imageChannel, sc)
	go drawFrame(imageChannel)
	if err := ebiten.Run(update, 1280, 720, 1, "IMAGE FROM WEBCAM"); err != nil {
		log.Fatal(err)
	}
}

func drawFrame(ic chan []byte) {
	for i := range ic {
		go frame(i, canvasImg)
	}
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.DrawImage(canvasImg, nil)
	ebiten.SetRunnableInBackground(true)
	return nil
}

func frame(b []byte, canvas *ebiten.Image) {
	newImage, err := toImage(b)
	if err != nil {
		log.Fatal(err)
	}
	camimg, err = ebiten.NewImageFromImage(newImage, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	canvas.DrawImage(camimg, nil)
}

func toImage(b []byte) (image.Image, error) {
	width := screenWidth
	height := screenHeight
	step := 3840
	data := b
	channels := 3

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	c := color.RGBA{
		R: uint8(0),
		G: uint8(0),
		B: uint8(0),
		A: uint8(255),
	}

	for y := 0; y < height; y++ {
		for x := 0; x < step; x = x + channels {
			c.B = uint8(data[y*step+x])
			c.G = uint8(data[y*step+x+1])
			c.R = uint8(data[y*step+x+2])
			if channels == 4 {
				c.A = uint8(data[y*step+x+3])
			}
			img.SetRGBA(int(x/channels), y, c)
		}
	}

	return img, nil
}
