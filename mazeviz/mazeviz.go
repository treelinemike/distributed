package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	for cols := 0; cols < 16; cols++ {
		for rows := 0; rows < 16; rows++ {
			cell := ebiten.NewImage(60, 60)
			cell.Fill(color.RGBA{0xc0, 0x00, 0xc0, 0xff})
			drawopts := new(ebiten.DrawImageOptions)
			drawopts.GeoM.Translate(float64(2+62*cols), float64(2+62*rows))
			screen.DrawImage(cell, drawopts)
		}
	}
	scale := ebiten.Monitor().DeviceScaleFactor()
	msg := fmt.Sprintf("Device Scale Ratio: %0.2f", scale)
	ebitenutil.DebugPrint(screen, msg)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	s := ebiten.Monitor().DeviceScaleFactor()
	return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)
}

func main() {
	ebiten.SetWindowSize(994, 994)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
