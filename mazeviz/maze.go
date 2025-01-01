package mazeviz

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Maze struct {
	cells []Cell
	//lines []Line
}

func pixdist(x0, y0, x1, y1 float64) float64 {
	return math.Sqrt(math.Pow((x1-x0), 2.0) + math.Pow((y1-y0), 2))
}

func (m *Maze) Draw(screen *ebiten.Image) {
	for _, c := range m.cells {
		c.Draw(screen)
	}
}

func NewMaze(p Params) (*Maze, error) {
	m := new(Maze)
	for cols := 0; cols < p.N; cols++ {
		for rows := 0; rows < p.M; rows++ {
			c := new(Cell)
			c.Size = p.CSZ
			c.X = p.CSP + p.CSK*cols
			c.Y = p.CSP + p.CSK*rows
			c.CX = float64(p.CSP) + float64(p.CSK)*(float64(cols)+0.5)
			c.CY = float64(p.CSP) + float64(p.CSK)*(float64(rows)+0.5)
			c.Color = color.RGBA{0x50, 0x50, 0x50, 0xff}
			m.cells = append(m.cells, *c)
		}
	}
	return m, nil
}
