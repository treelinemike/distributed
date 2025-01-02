package mazeviz

import (
	"errors"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Params struct {
	M   int
	N   int
	CSZ int
	CSP int
	CSK int
	WW  int
	WH  int
}

func (p *Params) Setparams(m int, n int) {
	p.M = m
	p.N = n

	numcells := max(m, n)
	p.CSZ = int(math.Round(708.0 / (1.1*float64(numcells) + 0.1)))
	p.CSP = int(math.Floor(((708.0 - float64(numcells*p.CSZ)) / float64(numcells+1))))
	p.CSK = p.CSZ + p.CSP
	p.WW = p.CSP + n*p.CSK
	p.WH = p.CSP + m*p.CSK

	/*
		fmt.Printf("csz: %d\n", p.CSZ)
		fmt.Printf("csp: %d\n", p.CSP)
		fmt.Printf("csk: %d\n", p.CSK)
		fmt.Printf("ww: %d\n", p.WW)
		fmt.Printf("wh: %d\n", p.WH)
	*/
}

type Game struct {
	image *ebiten.Image
	maze  *Maze
}

func NewGame(p Params) (*Game, error) {
	g := new(Game)
	g.image = new(ebiten.Image)
	m, err := NewMaze(p)
	if err != nil {
		return g, errors.New("couldn't initialize maze")
	}
	g.maze = m
	return g, nil
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		min_idx := 0
		min_dist := 1e6
		for i, c := range g.maze.cells {
			dist := pixdist(float64(x), float64(y), c.CX, c.CY)
			if dist < min_dist {
				min_dist = dist
				min_idx = i
			}
		}
		if g.maze.cells[min_idx].Color.G == 0x50 {
			g.maze.cells[min_idx].Color = color.RGBA{0x50, 0x00, 0x00, 0xff}
		} else {
			g.maze.cells[min_idx].Color = color.RGBA{0x50, 0x50, 0x50, 0xff}
		}
		fmt.Printf("Click at (%d,%d) -> cell %d\n", x, y, min_idx)

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.maze.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	s := ebiten.Monitor().DeviceScaleFactor()
	return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)
}
