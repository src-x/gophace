package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten"
	stan "github.com/nats-io/stan.go"
	"gocv.io/x/gocv"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

var camimg *ebiten.Image
var canvasImg *ebiten.Image
var wg sync.WaitGroup

var newMat gocv.Mat
var newImage image.Image

func init() {
	canvasImg, _ = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
	canvasImg.Fill(color.White)
}
func subscribe(c chan []byte, sc stan.Conn) {
	_, err := sc.Subscribe("foo", func(m *stan.Msg) {
		fmt.Print(".")
		c <- m.Data
	})
	if err != nil {
		log.Fatalln(err)
	}
	wg.Done()
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
	wg.Add(1)
	go subscribe(imageChannel, sc)
	wg.Wait()
	go drawFrame(imageChannel)
	if err := ebiten.Run(update, 1280, 720, 1, "IMAGE FROM WEBCAM"); err != nil {
		log.Fatal(err)
	}
}

func drawFrame(ic chan []byte) {
	for i := range ic {
		fmt.Println(i)
		frame(i, canvasImg)
	}
}
func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.Fill(color.Black)

	screen.DrawImage(canvasImg, nil)

	return nil

}

func frame(b []byte, canvas *ebiten.Image) {
	if camimg != nil {
		camimg.Clear()
	}
	newMat, err := gocv.NewMatFromBytes(720, 1280, 16, b)
	if err != nil {
		log.Fatal(err)
	}
	newImage, err := newMat.ToImage()
	if err != nil {
		log.Fatal(err)
	}
	camimg, err = ebiten.NewImageFromImage(newImage, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	canvas.DrawImage(camimg, nil)
}