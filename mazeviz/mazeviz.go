package main

import (
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
			cell := ebiten.NewImage(50, 50)
			cell.Fill(color.RGBA{0xc0, 0x00, 0xc0, 0xff})
			drawopts := new(ebiten.DrawImageOptions)
			drawopts.GeoM.Translate(float64(2+52*cols), float64(2+52*rows))
			screen.DrawImage(cell, drawopts)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	s := ebiten.Monitor().DeviceScaleFactor()
	return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)
}

func main() {
	ebiten.SetWindowSize(834, 834)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
