package mazeviz

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var params struct {
	m   int
	n   int
	csz int
	csp int
	csk int
	ww  int
	wh  int
}

func Setparams(m int, n int) {
	params.m = m
	params.n = n

	numcells := max(m, n)
	params.csz = int(math.Round(708.0 / (1.1*float64(numcells) + 0.1)))
	params.csp = int(math.Floor(((708.0 - float64(numcells*params.csz)) / float64(numcells+1))))
	params.csk = params.csz + params.csp
	params.ww = params.csp + n*params.csk
	params.wh = params.csp + m*params.csk

	fmt.Printf("csz: %d\n", params.csz)
	fmt.Printf("csp: %d\n", params.csp)
	fmt.Printf("csk: %d\n", params.csk)
}

func Start(windowtitle string) {
	ebiten.SetWindowSize(params.ww, params.wh)
	ebiten.SetWindowTitle(windowtitle)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//ebitenutil.DebugPrint(screen, "Hello, World!")
	for cols := 0; cols < params.n; cols++ {
		for rows := 0; rows < params.m; rows++ {
			cell := ebiten.NewImage(params.csz, params.csz)
			cell.Fill(color.RGBA{0x50, 0x50, 0x50, 0xff})
			drawopts := new(ebiten.DrawImageOptions)
			drawopts.GeoM.Translate(float64(params.csp+params.csk*cols), float64(params.csp+params.csk*rows))
			screen.DrawImage(cell, drawopts)
		}
	}
	//scale := ebiten.Monitor().DeviceScaleFactor()

	//msg := fmt.Sprintf("Device Scale Ratio: %0.2f", scale)
	//ebitenutil.DebugPrint(screen, msg)

	var path vector.Path
	path.MoveTo(0, 0)
	path.LineTo(100, 100)
	//path.AppendVerticesAndIndicesForStroke()

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	s := ebiten.Monitor().DeviceScaleFactor()
	return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)
}
