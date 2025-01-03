package mazeviz

import (
	"errors"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Params struct {
	M        int
	N        int
	CSZ      int
	CSP      int
	CSK      int
	WW       int
	WH       int
	Walls    []int32
	Cellvals []int32
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

		// search lines first
		for i, l := range g.maze.lines {
			if (float32(x) >= l.Xmin) && (float32(x) <= l.Xmax) && (float32(y) >= l.Ymin) && (float32(y) <= l.Ymax) {
				if l.Type == W_none {
					g.maze.lines[i].Type = W_true
				} else {
					g.maze.lines[i].Type = W_none
				}
				break
			}
		}

		// now search cells
		for i, c := range g.maze.cells {
			if (float32(x) > c.Xmin) && (float32(x) < c.Xmax) && (float32(y) > c.Ymin) && (float32(y) < c.Ymax) {
				if c.Type == C_none {
					g.maze.cells[i].Type = C_goal
				} else {
					g.maze.cells[i].Type = C_none
				}
				break
			}
		}

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
